package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		n := len(name)
		if n > 42 {
			n = 42
		}

		for i := 0; i < n; i++ {
			person.name[i] = name[i]
		}

		if n < 42 {
			person.name[len(name)] = 0
		}
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = int32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.stats = (person.stats &^ (0x3FF << 21)) | (int32(mana) << 21)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.health = uint16(health)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.stats = (person.stats &^ (0xF << 17)) | (int32(respect) << 17)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.stats = (person.stats &^ (0xF << 13)) | (int32(strength) << 13)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.stats = (person.stats &^ (0xF << 9)) | (int32(experience) << 9)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.stats = (person.stats &^ (0xF << 5)) | (int32(level) << 5)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.stats = (person.stats &^ (0x1 << 4)) | (int32(1) << 4)
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.stats = (person.stats &^ (0x1 << 3)) | (int32(1) << 3)
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.stats = (person.stats &^ (0x1 << 2)) | (int32(1) << 2)
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.stats = (person.stats &^ (0x3 << 0)) | (int32(personType) << 0)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	name          [42]byte
	health        uint16
	x, y, z, gold int32
	stats         int32
}

func NewGamePerson(options ...Option) GamePerson {
	gp := new(GamePerson)

	for _, o := range options {
		o(gp)
	}

	return *gp
}

func (p *GamePerson) Name() string {
	return string(p.name[:])
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	return int((p.stats >> 21) & 0x3FF)
}

func (p *GamePerson) Health() int {
	return int(p.health)
}

func (p *GamePerson) Respect() int {
	return int((p.stats >> 17) & 0xF)
}

func (p *GamePerson) Strength() int {
	return int((p.stats >> 13) & 0xF)
}

func (p *GamePerson) Experience() int {
	return int((p.stats >> 9) & 0xF)
}

func (p *GamePerson) Level() int {
	return int((p.stats >> 5) & 0xF)
}

func (p *GamePerson) HasHouse() bool {
	return ((p.stats >> 4) & 0x1) != 0
}

func (p *GamePerson) HasGun() bool {
	return ((p.stats >> 3) & 0x1) != 0
}

func (p *GamePerson) HasFamilty() bool {
	return ((p.stats >> 2) & 0x1) != 0
}

func (p *GamePerson) Type() int {
	return int((p.stats) & 0x3)
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
