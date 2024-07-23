package maths

import "math"

type Number interface {
	int | float64
}

type Vector[T Number] struct {
	X, Y, Z T
}

func (v *Vector[T]) Clone() Vector[T] {
	return Vector[T]{v.X, v.Y, v.Z}
}

func (v *Vector[T]) Sum(w Vector[T]) {
	v.X += w.X
	v.Y += w.Y
	v.Z += w.Z
}

func (v *Vector[T]) Sub(w Vector[T]) {
	v.X -= w.X
	v.Y -= w.Y
	v.Z -= w.Z
}

func (v *Vector[T]) Mul(w Vector[T]) {
	v.X *= w.X
	v.Y *= w.Y
	v.Z *= w.Z
}

func (v *Vector[T]) Div(w Vector[T]) {
	v.X /= w.X
	v.Y /= w.Y
	v.Z /= w.Z
}

func (v *Vector[T]) Scale(s T) {
	v.X *= s
	v.Y *= s
	v.Z *= s
}

func (v *Vector[T]) Heading() float64 {
	return math.Atan2(float64(v.Y), float64(v.X)) * (180 / math.Pi)
}

func (v *Vector[T]) Clamp(min, max Vector[T]) {
	v.X = clamp(v.X, min.X, max.X)
	v.Y = clamp(v.Y, min.Y, max.Y)
	v.Z = clamp(v.Z, min.Z, max.Z)
}

func clamp[T Number](value, min, max T) T {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
