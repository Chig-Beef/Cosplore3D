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

type audioStream interface {
	io.ReadSeeker
	Length() int64
}

func (g *Game) load_audio() error {
	g.audio = make(map[string]*AudioPlayer)

	if err := g.load_track("Ankaran", "ankaran"); err != nil {
		return err
	}
	if err := g.load_track("Enikoko", "enikoko"); err != nil {
		return err
	}
	if err := g.load_track("enemyDeath", "enemyDeath"); err != nil {
		return err
	}
	if err := g.load_track("enemyHurt", "enemyHurt"); err != nil {
		return err
	}
	if err := g.load_track("pickup", "pickup"); err != nil {
		return err
	}
	if err := g.load_track("playerHurt", "playerHurt"); err != nil {
		return err
	}
	if err := g.load_track("shoot", "shoot"); err != nil {
		return err
	}
	if err := g.load_track("trigger", "trigger"); err != nil {
		return err
	}

	return nil
}

func (g *Game) load_track(fName, mName string) error {
	var s audioStream

	audio, err := os.ReadFile("assets/audio/" + fName + ".ogg")
	if err != nil {
		return err
	}

	s, err = vorbis.DecodeWithoutResampling(bytes.NewReader(audio))
	if err != nil {
		log.Fatal("failed to load audio")
	}

	ap, err := NewPlayer(g, 0, s)
	if err != nil {
		return err
	}

	g.audio[mName] = ap

	return nil
}

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

func NewPlayer(g *Game, musicType musicType, stream audioStream) (*AudioPlayer, error) {
	const bytesPerSample = 4

	p, err := g.ctx.NewPlayer(stream)
	if err != nil {
		return nil, err
	}

	player := &AudioPlayer{
		game:         g,
		audioContext: g.ctx,
		audioPlayer:  p,
		total:        time.Second * time.Duration(stream.Length()) / bytesPerSample / time.Duration(sampleRate),
		volume128:    128,
		musicType:    musicType,
	}
	if player.total == 0 {
		player.total = 1
	}

	return player, nil
}

func (ap *AudioPlayer) Close() error {
	return ap.audioPlayer.Close()

}

func (ap *AudioPlayer) update(g *Game) error {
	if ap.audioPlayer.IsPlaying() {
		ap.current = ap.audioPlayer.Position()
	} else {
		err := ap.audioPlayer.Rewind()
		if err != nil {
			return err
		}
		ap.audioPlayer.Play()
	}
	return nil
}

func (g *Game) play_audio(key string) {
	m, ok := g.audio[g.curAudio]
	if ok {
		if m.audioPlayer.IsPlaying() {
			m.audioPlayer.Pause()
		}
	}

	m = g.audio[key]

	g.curAudio = key
	g.musicPlayer = m

	err := g.musicPlayer.audioPlayer.Rewind()
	if err != nil {
		log.Fatal(err)
	}
	g.musicPlayer.audioPlayer.Play()
}

func (g *Game) init_audio() error {
	audioContext := audio.NewContext(sampleRate)

	g.musicPlayerCh = make(chan *AudioPlayer)
	g.errCh = make(chan error)

	g.ctx = audioContext

	err := g.load_audio()

	return err
}
