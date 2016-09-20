// Copyright 2012 Daniel Connelly.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtreego

import (
	"math"
	"testing"
)

const EPS = 0.000000001

func TestDist(t *testing.T) {
	p := Point{1, 2}
	q := Point{4, 5}
	dist := math.Sqrt(18)
	if d := p.dist(q); d != dist {
		t.Errorf("dist(%v, %v) = %v; expected %v", p, q, d, dist)
	}
}

func TestNewRect(t *testing.T) {
	p := Point{1.0, -2.5}
	q := Point{3.5, 5.5}
	lengths := [Dim]float64{2.5, 8.0}

	rect, err := NewRect(p, lengths)
	if err != nil {
		t.Errorf("Error on NewRect(%v, %v): %v", p, lengths, err)
	}
	if d := p.dist(rect.p); d > EPS {
		t.Errorf("Expected p == rect.p")
	}
	if d := q.dist(rect.q); d > EPS {
		t.Errorf("Expected q == rect.q")
	}
}

func TestNewRectDistError(t *testing.T) {
	p := Point{1.0, -2.5}
	lengths := [Dim]float64{2.5, -8.0}
	_, err := NewRect(p, lengths)
	if _, ok := err.(DistError); !ok {
		t.Errorf("Expected distError on NewRect(%v, %v)", p, lengths)
	}
}

func TestRectSize(t *testing.T) {
	p := Point{1.0, -2.5}
	lengths := [Dim]float64{2.5, 8.0}
	rect, _ := NewRect(p, lengths)
	size := lengths[0] * lengths[1]
	actual := rect.size()
	if size != actual {
		t.Errorf("Expected %v.size() == %v, got %v", rect, size, actual)
	}
}

func TestRectMargin(t *testing.T) {
	p := Point{1.0, -2.5}
	lengths := [Dim]float64{2.5, 8.0}
	rect, _ := NewRect(p, lengths)
	size := 4*2.5 + 4*8.0
	actual := rect.margin()
	if size != actual {
		t.Errorf("Expected %v.margin() == %v, got %v", rect, size, actual)
	}
}

func TestContainsPoint(t *testing.T) {
	p := Point{3.7, -2.4}
	lengths := [Dim]float64{6.2, 1.1}
	rect, _ := NewRect(p, lengths)

	q := Point{4.5, -1.7}
	if yes := rect.containsPoint(q); !yes {
		t.Errorf("Expected %v contains %v", rect, q)
	}
}

func TestDoesNotContainPoint(t *testing.T) {
	p := Point{3.7, -2.4}
	lengths := [Dim]float64{6.2, 1.1}
	rect, _ := NewRect(p, lengths)

	q := Point{4.5, -2.7}
	if yes := rect.containsPoint(q); yes {
		t.Errorf("Expected %v doesn't contain %v", rect, q)
	}

	q = Point{3.6, -2.0}
	if yes := rect.containsPoint(q); yes {
		t.Errorf("Expected %v doesn't contain %v", rect, q)
	}
}

func TestContainsRect(t *testing.T) {
	p := Point{3.7, -2.4}
	lengths1 := [Dim]float64{6.2, 1.1}
	rect1, _ := NewRect(p, lengths1)

	q := Point{4.1, -1.9}
	lengths2 := [Dim]float64{3.2, 0.6}
	rect2, _ := NewRect(q, lengths2)

	if yes := rect1.containsRect(&rect2); !yes {
		t.Errorf("Expected %v.containsRect(%v", rect1, rect2)
	}
}

func TestDoesNotContainRectOverlaps(t *testing.T) {
	p := Point{3.7, -2.4}
	lengths1 := [Dim]float64{6.2, 1.1}
	rect1, _ := NewRect(p, lengths1)

	q := Point{4.1, -1.9}
	lengths2 := [Dim]float64{3.2, 1.4}
	rect2, _ := NewRect(q, lengths2)

	if yes := rect1.containsRect(&rect2); yes {
		t.Errorf("Expected %v doesn't contain %v", rect1, rect2)
	}
}

func TestDoesNotContainRectDisjoint(t *testing.T) {
	p := Point{3.7, -2.4}
	lengths1 := [Dim]float64{6.2, 1.1}
	rect1, _ := NewRect(p, lengths1)

	q := Point{1.2, -19.6}
	lengths2 := [Dim]float64{2.2, 5.9}
	rect2, _ := NewRect(q, lengths2)

	if yes := rect1.containsRect(&rect2); yes {
		t.Errorf("Expected %v doesn't contain %v", rect1, rect2)
	}
}

func TestNoIntersection(t *testing.T) {
	p := Point{1, 2}
	lengths1 := [Dim]float64{1, 1}
	rect1, _ := NewRect(p, lengths1)

	q := Point{-1, -2}
	lengths2 := [Dim]float64{2.5, 3}
	rect2, _ := NewRect(q, lengths2)

	// rect1 and rect2 fail to overlap in just one dimension (second)

	if intersect(&rect1, &rect2) {
		t.Errorf("Expected intersect(%v, %v) == false", rect1, rect2)
	}
}

func TestNoIntersectionJustTouches(t *testing.T) {
	p := Point{1, 2}
	lengths1 := [Dim]float64{1, 1}
	rect1, _ := NewRect(p, lengths1)

	q := Point{-1, -2}
	lengths2 := [Dim]float64{2.5, 4}
	rect2, _ := NewRect(q, lengths2)

	// rect1 and rect2 fail to overlap in just one dimension (second)

	if intersect(&rect1, &rect2) {
		t.Errorf("Expected intersect(%v, %v) == nil", rect1, rect2)
	}
}

func TestContainmentIntersection(t *testing.T) {
	p := Point{1, 2}
	lengths1 := [Dim]float64{1, 1}
	rect1, _ := NewRect(p, lengths1)

	q := Point{1, 2.2}
	lengths2 := [Dim]float64{0.5, 0.5}
	rect2, _ := NewRect(q, lengths2)

	r := Point{1, 2.2}
	s := Point{1.5, 2.7}

	if !intersect(&rect1, &rect2) {
		t.Errorf("intersect(%v, %v) != %v, %v", rect1, rect2, r, s)
	}
}

func TestOverlapIntersection(t *testing.T) {
	p := Point{1, 2}
	lengths1 := [Dim]float64{1, 2.5}
	rect1, _ := NewRect(p, lengths1)

	q := Point{1, 4}
	lengths2 := [Dim]float64{3, 2}
	rect2, _ := NewRect(q, lengths2)

	r := Point{1, 4}
	s := Point{2, 4.5}

	if !intersect(&rect1, &rect2) {
		t.Errorf("intersect(%v, %v) != %v, %v", rect1, rect2, r, s)
	}
}

func TestToRect(t *testing.T) {
	x := Point{3.7, -2.4}
	tol := 0.05
	rect := x.ToRect(tol)

	p := Point{3.65, -2.45}
	q := Point{3.75, -2.35}
	d1 := p.dist(rect.p)
	d2 := q.dist(rect.q)
	if d1 > EPS || d2 > EPS {
		t.Errorf("Expected %v.ToRect(%v) == %v, %v, got %v", x, tol, p, q, rect)
	}
}

func TestBoundingBox(t *testing.T) {
	p := Point{3.7, -2.4}
	lengths1 := [Dim]float64{1, 15}
	rect1, _ := NewRect(p, lengths1)

	q := Point{-6.5, 4.7}
	lengths2 := [Dim]float64{4, 5}
	rect2, _ := NewRect(q, lengths2)

	r := Point{-6.5, -2.4}
	s := Point{4.7, 12.6}

	bb := boundingBox(&rect1, &rect2)
	d1 := r.dist(bb.p)
	d2 := s.dist(bb.q)
	if d1 > EPS || d2 > EPS {
		t.Errorf("boundingBox(%v, %v) != %v, %v, got %v", rect1, rect2, r, s, bb)
	}
}

func TestBoundingBoxContains(t *testing.T) {
	p := Point{3.7, -2.4}
	lengths1 := [Dim]float64{1, 15}
	rect1, _ := NewRect(p, lengths1)

	q := Point{4.0, 0.0}
	lengths2 := [Dim]float64{0.56, 6.222222}
	rect2, _ := NewRect(q, lengths2)

	bb := boundingBox(&rect1, &rect2)
	d1 := rect1.p.dist(bb.p)
	d2 := rect1.q.dist(bb.q)
	if d1 > EPS || d2 > EPS {
		t.Errorf("boundingBox(%v, %v) != %v, got %v", rect1, rect2, rect1, bb)
	}
}

func TestMinDistZero(t *testing.T) {
	p := Point{1, 2}
	r := p.ToRect(1)
	if d := p.minDist(r); d > EPS {
		t.Errorf("Expected %v.minDist(%v) == 0, got %v", p, r, d)
	}
}

func TestMinDistPositive(t *testing.T) {
	p := Point{1, 2}
	r := Rect{Point{-1, -4}, Point{2, -2}}
	expected := float64((-2 - 2) * (-2 - 2))
	if d := p.minDist(&r); math.Abs(d-expected) > EPS {
		t.Errorf("Expected %v.minDist(%v) == %v, got %v", p, r, expected, d)
	}
}

func TestMinMaxdist(t *testing.T) {
	p := Point{-3, -2}
	r := Rect{Point{0, 0}, Point{1, 2}}

	// furthest points from p on the faces closest to p in each dimension
	q1 := Point{0, 2}
	q2 := Point{1, 0}
	q3 := Point{1, 2}

	// find the closest distance from p to one of these furthest points
	d1 := p.dist(q1)
	d2 := p.dist(q2)
	d3 := p.dist(q3)
	expected := math.Min(d1*d1, math.Min(d2*d2, d3*d3))

	if d := p.minMaxDist(&r); math.Abs(d-expected) > EPS {
		t.Errorf("Expected %v.minMaxDist(%v) == %v, got %v", p, r, expected, d)
	}
}
