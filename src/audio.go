package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"slices"
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
	g.musicPlayerCh = make(map[string]chan *AudioPlayer)
	g.errCh = make(chan error)
	g.audio = make(map[string]*AudioPlayer)
	g.soundEffects = make(map[string]*AudioPlayer)
	g.musicPlayer = make(map[string]*AudioPlayer)

	if err := g.load_track("Ankaran", "ankaran"); err != nil {
		return err
	}
	if err := g.load_track("Cosplorer", "cosplorer"); err != nil {
		return err
	}
	if err := g.load_track("Enikoko", "enikoko"); err != nil {
		return err
	}
	if err := g.load_effect("enemyDeath", "enemyDeath"); err != nil {
		return err
	}
	if err := g.load_effect("enemyHurt", "enemyHurt"); err != nil {
		return err
	}
	if err := g.load_effect("pickup", "pickup"); err != nil {
		return err
	}
	if err := g.load_effect("playerHurt", "playerHurt"); err != nil {
		return err
	}
	if err := g.load_effect("shoot", "shoot"); err != nil {
		return err
	}
	if err := g.load_effect("trigger", "trigger"); err != nil {
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

func (g *Game) load_effect(fName, mName string) error {
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

	g.soundEffects[mName] = ap

	return nil
}

func (g *Game) update_audio() error {
	if _, ok := g.musicPlayerCh[g.curAudio]; ok {
		select {
		case p := <-g.musicPlayerCh[g.curAudio]:
			g.musicPlayer[g.curAudio] = p
		case err := <-g.errCh:
			return err
		default:
		}

		if g.musicPlayer[g.curAudio] != nil {
			if err := g.musicPlayer[g.curAudio].update(g); err != nil {
				return err
			}
		}
	}

	for _, sfx := range g.curSoundEffects {
		if _, ok := g.musicPlayerCh[sfx]; ok {
			select {
			case p := <-g.musicPlayerCh[sfx]:
				g.musicPlayer[sfx] = p
			case err := <-g.errCh:
				return err
			default:
			}

			if g.musicPlayer[sfx] != nil {
				if err := g.musicPlayer[sfx].update_as_effect(g); err != nil {
					return err
				}
			}
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

func (ap *AudioPlayer) update_as_effect(g *Game) error {
	if ap.audioPlayer.IsPlaying() {
		ap.current = ap.audioPlayer.Position()
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
	g.musicPlayer[key] = m

	err := g.musicPlayer[key].audioPlayer.Rewind()
	if err != nil {
		log.Fatal(err)
	}
	g.musicPlayer[key].audioPlayer.Play()
}

func (g *Game) play_effect(key string) {
	m, ok := g.soundEffects[g.curAudio]
	if ok {
		if m.audioPlayer.IsPlaying() {
			m.audioPlayer.Pause()
		}
	}

	m = g.soundEffects[key]

	if !slices.Contains(g.curSoundEffects, key) {
		g.curSoundEffects = append(g.curSoundEffects, key)
	}
	g.musicPlayer[key] = m

	err := g.musicPlayer[key].audioPlayer.Rewind()
	if err != nil {
		log.Fatal(err)
	}
	g.musicPlayer[key].audioPlayer.Play()
}

func (g *Game) init_audio() error {
	audioContext := audio.NewContext(sampleRate)
	g.ctx = audioContext

	err := g.load_audio()

	return err
}
