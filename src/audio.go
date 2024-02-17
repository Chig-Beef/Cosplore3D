package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

type AudioPlayer struct {
	game         *Game
	audioContext *audio.Context
	audioPlayer  *audio.Player
	current      time.Duration
	total        time.Duration
	volume128    int
	musicType    musicType
}

type musicType int

const sampleRate int = 48_000

func (g *Game) update_audio() error {
	select {
	case p := <-g.musicPlayerCh:
		g.musicPlayer = p
	case err := <-g.errCh:
		return err
	default:
	}

	if g.musicPlayer != nil {
		if err := g.musicPlayer.update(g); err != nil {
			return err
		}
	}

	return nil
}

func NewPlayer(g *Game, audioContext *audio.Context, musicType musicType) (*AudioPlayer, error) {
	type audioStream interface {
		io.ReadSeeker
		Length() int64
	}

	const bytesPerSample = 4

	var s audioStream

	ankaranAudio, err := os.ReadFile("assets/audio/Ankaran.ogg")
	if err != nil {
		return nil, err
	}

	s, err = vorbis.DecodeWithoutResampling(bytes.NewReader(ankaranAudio))
	if err != nil {
		log.Fatal("failed to load audio")
	}

	p, err := audioContext.NewPlayer(s)
	if err != nil {
		return nil, err
	}

	player := &AudioPlayer{
		game:         g,
		audioContext: audioContext,
		audioPlayer:  p,
		total:        time.Second * time.Duration(s.Length()) / bytesPerSample / time.Duration(sampleRate),
		volume128:    128,
		musicType:    musicType,
	}
	if player.total == 0 {
		player.total = 1
	}

	player.audioPlayer.Play()
	return player, nil
}

func (ap *AudioPlayer) Close() error {
	return ap.audioPlayer.Close()

}

func (ap *AudioPlayer) update(g *Game) error {
	if ap.audioPlayer.IsPlaying() {
		ap.current = ap.audioPlayer.Position()
	} else {
		ap.playBackground(g)
	}
	return nil
}

func (ap *AudioPlayer) playBackground(g *Game) {
	m, err := NewPlayer(g, g.musicPlayer.audioContext, 0)
	if err != nil {
		log.Fatal("error with sound")
	}

	g.musicPlayer = m
}

func (g *Game) init_audio() {
	audioContext := audio.NewContext(sampleRate)
	m, err := NewPlayer(g, audioContext, 0)
	if err != nil {
		log.Fatal("error with sound")
	}
	g.musicPlayer = m

	g.musicPlayerCh = make(chan *AudioPlayer)
	g.errCh = make(chan error)
}
