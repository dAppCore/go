package core_test

import (
	. "dappco.re/go"
)

func TestMath_Min_Good(t *T) {
	AssertEqual(t, 3, Min(3, 7))
	AssertEqual(t, 3, Min(7, 3)) // order independent
	AssertEqual(t, -5, Min(-5, -2))
	AssertEqual(t, -5, Min(-5, 5))
	AssertEqual(t, 1.5, Min(2.5, 1.5)) // float64
	// Go's builtin min propagates NaN: any NaN operand yields NaN.
	AssertTrue(t, IsNaN(Min(NaN(), 1.0)))
	AssertTrue(t, IsNaN(Min(1.0, NaN())))
}

func TestMath_Min_Bad(t *T) {
	AssertEqual(t, 3, Min(3, 3)) // equal ints
	AssertEqual(t, 0, Min(0, 0))
	AssertEqual(t, -4, Min(-4, -4))
	AssertEqual(t, "x", Min("x", "x")) // equal strings
	AssertEqual(t, 2.5, Min(2.5, 2.5)) // equal floats
}

func TestMath_Min_Ugly(t *T) {
	AssertEqual(t, "alpha", Min("beta", "alpha"))
	AssertEqual(t, "alpha", Min("alpha", "beta")) // order independent
	AssertEqual(t, "", Min("", "a"))              // empty string sorts first
	AssertEqual(t, "Z", Min("Z", "a"))            // uppercase < lowercase (ASCII)
	AssertEqual(t, "apple", Min("apple", "apply"))
}

func TestMath_Max_Good(t *T) {
	AssertEqual(t, 7, Max(3, 7))
	AssertEqual(t, 7, Max(7, 3)) // order independent
	AssertEqual(t, -2, Max(-5, -2))
	AssertEqual(t, 5, Max(-5, 5))
	AssertEqual(t, 2.5, Max(2.5, 1.5)) // float64
	// Go's builtin max propagates NaN: any NaN operand yields NaN.
	AssertTrue(t, IsNaN(Max(NaN(), 1.0)))
	AssertTrue(t, IsNaN(Max(1.0, NaN())))
}

func TestMath_Max_Bad(t *T) {
	AssertEqual(t, 3, Max(3, 3)) // equal ints
	AssertEqual(t, 0, Max(0, 0))
	AssertEqual(t, -4, Max(-4, -4))
	AssertEqual(t, "x", Max("x", "x")) // equal strings
	AssertEqual(t, 2.5, Max(2.5, 2.5)) // equal floats
}

func TestMath_Max_Ugly(t *T) {
	AssertEqual(t, "beta", Max("beta", "alpha"))
	AssertEqual(t, "beta", Max("alpha", "beta")) // order independent
	AssertEqual(t, "a", Max("", "a"))            // non-empty sorts after empty
	AssertEqual(t, "a", Max("Z", "a"))           // lowercase > uppercase (ASCII)
	AssertEqual(t, "apply", Max("apple", "apply"))
}

func TestMath_Abs_Good(t *T) {
	AssertEqual(t, 42, Abs(-42))
	AssertEqual(t, 1, Abs(-1))
	AssertEqual(t, 100, Abs(-100))
	AssertEqual(t, int64(9), Abs(int64(-9))) // wider signed type
	AssertEqual(t, 2.5, Abs(-2.5))           // float64
}

func TestMath_Abs_Bad(t *T) {
	AssertEqual(t, 0, Abs(0))
	AssertEqual(t, 42, Abs(42))   // already positive, unchanged
	AssertEqual(t, 0.0, Abs(0.0)) // float zero
	AssertEqual(t, 7.5, Abs(7.5)) // positive float unchanged
}

func TestMath_Abs_Ugly(t *T) {
	AssertEqual(t, float32(3.5), Abs(float32(-3.5)))
	AssertEqual(t, float32(3.5), Abs(float32(3.5)))
	// Two's-complement overflow edge: -MinInt8 wraps back to MinInt8.
	AssertEqual(t, int8(-128), Abs(int8(-128)))
	AssertEqual(t, int8(127), Abs(int8(-127))) // -127 is representable
}

func TestMath_Clamp_Good(t *T) {
	AssertEqual(t, 5, Clamp(5, 0, 10))        // inside the range, unchanged
	AssertEqual(t, 3.5, Clamp(3.5, 1.0, 9.0)) // float64
	AssertEqual(t, int64(8), Clamp(int64(8), int64(0), int64(16)))
}

func TestMath_Clamp_Bad(t *T) {
	AssertEqual(t, 0, Clamp(-4, 0, 10))  // below lo → lo
	AssertEqual(t, 10, Clamp(99, 0, 10)) // above hi → hi
	AssertEqual(t, 1.0, Clamp(-2.5, 1.0, 9.0))
}

func TestMath_Clamp_Ugly(t *T) {
	AssertEqual(t, 0, Clamp(0, 0, 10))   // lo boundary is inclusive
	AssertEqual(t, 10, Clamp(10, 0, 10)) // hi boundary is inclusive
	AssertEqual(t, 7, Clamp(3, 7, 7))    // lo==hi pins to the single value
	AssertEqual(t, 7, Clamp(99, 7, 7))
}

func TestMath_Sign_Good(t *T) {
	AssertEqual(t, 1, Sign(42))      // positive int → 1
	AssertEqual(t, 1.0, Sign(0.001)) // positive float → 1
	AssertEqual(t, int64(1), Sign(int64(9)))
}

func TestMath_Sign_Bad(t *T) {
	AssertEqual(t, -1, Sign(-42))      // negative int → -1
	AssertEqual(t, -1.0, Sign(-0.001)) // negative float → -1
	AssertEqual(t, int64(-1), Sign(int64(-9)))
}

func TestMath_Sign_Ugly(t *T) {
	AssertEqual(t, 0, Sign(0))       // zero → 0
	AssertEqual(t, 0.0, Sign(0.0))   // float zero → 0
	AssertEqual(t, 0.0, Sign(NaN())) // NaN → 0 (comparisons are false)
}

func TestMath_Pow_Good(t *T) {
	AssertEqual(t, 16.0, Pow(4, 2))
	AssertEqual(t, 1024.0, Pow(2, 10))
	AssertEqual(t, 8.0, Pow(2, 3))
	AssertEqual(t, 5.0, Pow(5, 1))   // exponent 1 is identity
	AssertEqual(t, -8.0, Pow(-2, 3)) // negative base, odd exponent
	AssertEqual(t, 4.0, Pow(-2, 2))  // negative base, even exponent
}

func TestMath_Pow_Bad(t *T) {
	AssertEqual(t, 1.0, Pow(4, 0))
	AssertEqual(t, 1.0, Pow(0, 0))   // 0^0 is defined as 1
	AssertEqual(t, 0.0, Pow(0, 5))   // 0 to a positive power
	AssertEqual(t, 1.0, Pow(1, 100)) // 1 to any power
	AssertEqual(t, 0.5, Pow(2, -1))  // negative exponent
}

func TestMath_Pow_Ugly(t *T) {
	AssertInDelta(t, 3.0, Pow(9, 0.5), 0.000001)
	AssertInDelta(t, 2.0, Pow(8, 1.0/3.0), 0.000001) // cube root
	AssertInDelta(t, 0.5, Pow(4, -0.5), 0.000001)    // reciprocal square root
	AssertTrue(t, IsNaN(Pow(-1, 0.5)))               // root of negative is NaN
	AssertTrue(t, IsNaN(Pow(-4, 0.5)))
}

func TestMath_Floor_Good(t *T) {
	AssertEqual(t, 3.0, Floor(3.7))
	AssertEqual(t, 3.0, Floor(3.2))
	AssertEqual(t, 0.0, Floor(0.9))
	AssertEqual(t, 5.0, Floor(5.999))
	AssertEqual(t, 2.0, Floor(2.0)) // exact integer unchanged
}

func TestMath_Floor_Bad(t *T) {
	AssertEqual(t, -4.0, Floor(-3.1))
	AssertEqual(t, -4.0, Floor(-3.9))
	AssertEqual(t, -1.0, Floor(-0.1))    // rounds toward negative infinity
	AssertEqual(t, -1.0, Floor(-0.0001)) // barely negative still floors to -1
	AssertEqual(t, -3.0, Floor(-3.0))    // exact negative integer unchanged
}

func TestMath_Floor_Ugly(t *T) {
	AssertEqual(t, 3.0, Floor(3))
	AssertEqual(t, 0.0, Floor(0))
	AssertEqual(t, -7.0, Floor(-7))
	AssertTrue(t, IsNaN(Floor(NaN()))) // NaN propagates
}

func TestMath_Ceil_Good(t *T) {
	AssertEqual(t, 4.0, Ceil(3.1))
	AssertEqual(t, 4.0, Ceil(3.9))
	AssertEqual(t, 1.0, Ceil(0.0001))
	AssertEqual(t, 6.0, Ceil(5.001))
	AssertEqual(t, 3.0, Ceil(3.0)) // exact integer unchanged
}

func TestMath_Ceil_Bad(t *T) {
	AssertEqual(t, -3.0, Ceil(-3.7))
	AssertEqual(t, -3.0, Ceil(-3.1)) // rounds toward positive infinity
	AssertEqual(t, -5.0, Ceil(-5.999))
	AssertEqual(t, -2.0, Ceil(-2.0)) // exact negative integer unchanged
}

func TestMath_Ceil_Ugly(t *T) {
	AssertEqual(t, 3.0, Ceil(3))
	AssertEqual(t, 0.0, Ceil(0))
	AssertEqual(t, -7.0, Ceil(-7))
	AssertTrue(t, IsNaN(Ceil(NaN()))) // NaN propagates
}

func TestMath_Round_Good(t *T) {
	AssertEqual(t, 4.0, Round(3.5))
	AssertEqual(t, 3.0, Round(2.5)) // half rounds away from zero
	AssertEqual(t, 1.0, Round(0.5))
	AssertEqual(t, 5.0, Round(4.5))
	AssertEqual(t, 2.0, Round(2.4)) // below half rounds down
}

func TestMath_Round_Bad(t *T) {
	AssertEqual(t, -4.0, Round(-3.5))
	AssertEqual(t, -1.0, Round(-0.5)) // negative half rounds away from zero
	AssertEqual(t, -3.0, Round(-2.5))
	AssertEqual(t, -2.0, Round(-2.4)) // above -half rounds toward zero
}

func TestMath_Round_Ugly(t *T) {
	AssertEqual(t, 3.0, Round(3.49))
	AssertEqual(t, 3.0, Round(3.0)) // exact integer unchanged
	AssertEqual(t, 0.0, Round(0.4)) // rounds to zero
	AssertEqual(t, 4.0, Round(3.51))
	AssertTrue(t, IsNaN(Round(NaN()))) // NaN propagates
}

func TestMath_NaN_Good(t *T) {
	n := NaN()

	AssertTrue(t, IsNaN(n))
	AssertTrue(t, IsNaN(n*2)) // NaN propagates through arithmetic
	AssertTrue(t, IsNaN(n+1))
	// All ordered comparisons against NaN are false.
	AssertFalse(t, n > 0)
	AssertFalse(t, n < 0)
	AssertFalse(t, n >= 0)
}

func TestMath_NaN_Bad(t *T) {
	AssertFalse(t, IsNaN(0.0))
	AssertFalse(t, IsNaN(1.5))
	AssertFalse(t, IsNaN(-1.5))
	AssertFalse(t, IsNaN(-0.0))
	AssertFalse(t, IsNaN(1e308))
	// +Inf is not NaN.
	AssertFalse(t, IsNaN(Pow(0, -1)))
}

func TestMath_NaN_Ugly(t *T) {
	n := NaN()
	m := n

	AssertFalse(t, n == m)
}

func TestMath_IsNaN_Good(t *T) {
	AssertTrue(t, IsNaN(NaN()))
	AssertTrue(t, IsNaN(Pow(-1, 0.5))) // root of a negative
	AssertTrue(t, IsNaN(NaN()+1))      // arithmetic preserves NaN
	inf := Pow(0, -1)
	AssertTrue(t, IsNaN(inf-inf)) // Inf - Inf is NaN
}

func TestMath_IsNaN_Bad(t *T) {
	AssertFalse(t, IsNaN(42))
	AssertFalse(t, IsNaN(0.0))
	AssertFalse(t, IsNaN(-7.5))
	AssertFalse(t, IsNaN(1e10))
	AssertFalse(t, IsNaN(Pow(0, -1))) // +Inf is not NaN
}

func TestMath_IsNaN_Ugly(t *T) {
	AssertTrue(t, IsNaN(Pow(-1, 0.5)))
	AssertTrue(t, IsNaN(Pow(-4, 0.5)))
	AssertTrue(t, IsNaN(NaN()*0)) // NaN * 0 is NaN, not 0
	inf := Pow(0, -1)
	AssertTrue(t, IsNaN(inf/inf)) // Inf / Inf is NaN
}

func TestMath_Compare_Good(t *T) {
	AssertEqual(t, -1, Compare(1, 2))
	AssertEqual(t, 1, Compare(2, 1))
	AssertEqual(t, -1, Compare(-5, -2))
	AssertEqual(t, 1, Compare(-2, -5))
	AssertEqual(t, -1, Compare(1.5, 2.5)) // floats
	AssertEqual(t, 1, Compare(2.5, 1.5))
	// cmp.Compare orders NaN below every non-NaN value.
	AssertEqual(t, -1, Compare(NaN(), 1.0))
	AssertEqual(t, 1, Compare(1.0, NaN()))
}

func TestMath_Compare_Bad(t *T) {
	AssertEqual(t, 0, Compare(7, 7))
	AssertEqual(t, 0, Compare(-3, -3))
	AssertEqual(t, 0, Compare(0, 0))
	AssertEqual(t, 0, Compare(2.5, 2.5))     // equal floats
	AssertEqual(t, 0, Compare("x", "x"))     // equal strings
	AssertEqual(t, 0, Compare(NaN(), NaN())) // cmp.Compare treats two NaNs as equal
}

func TestMath_Compare_Ugly(t *T) {
	AssertEqual(t, -1, Compare("alpha", "beta"))
	AssertEqual(t, 1, Compare("beta", "alpha"))
	AssertEqual(t, -1, Compare("", "a")) // empty string sorts first
	AssertEqual(t, 1, Compare("a", ""))
	AssertEqual(t, -1, Compare("Z", "a")) // uppercase < lowercase (ASCII)
	AssertEqual(t, -1, Compare("apple", "apply"))
	AssertEqual(t, 0, Compare("same", "same"))
}
