#!/usr/bin/python3
"""
This proof-of-concept in python aims to provide an optimized way (at least at the
best of my capacities) to search for maximum persistence in integers.

This is inspired by https://www.youtube.com/watch?v=Wim9WJeDTHQ

The first hard problem is to find a list of interesting integers to check as they
are not all 

The interesting combinations of length N are:
    2[6-9]{N-1},
    3[4,6-9][6-9]{N-2},
    [4,6-9][6-9]{N-1}

I'll exclude the 5 for now but it could be relevant to include it only if the
digits are exclusively odd numbers.

So the possible digits are [2-4,6-9].

Numbers are represented as an array of digits in base 7 (length of the list of
possible digits).
For example: 266699 is represented as [6,6,3,3,3,0]
                                       | | | | | |
                        (translate to) 9 9 6 6 6 2

The numbers generated are always in increasing order to easily avoid duplicate
permutations.
So a list of generated numbers of length 4 looks like this:
    2666 (actual starting number)
    2667
    2668
    2669
    2677
    2678
    2679
    2777
    2778
    2779
    2788
    2789
    2799
    2888
    ...
    2999
    3466
    3467
    3477
    ...
    3499
    3666
    3667
    ...
    9999

Having a number like 22xxx is the same as having 4xxx so this is avoided. Same
goes for a number like 23xxx => 6xxx, 24xxx => 8xxx or 33xxx => 9xxx.

"""
import sys
from math import ceil, log
from functools import reduce, lru_cache

#translation = [ 2, 3, 4, 5, 6, 7, 8, 9 ]
translation = [ 2, 3, 4, 6, 7, 8, 9 ]


def increment(arr, base):
    """Increments a 'number' recursively from the lowest to the highest digit.

    When the lowest digit attains the maximum allowed in the current base, it
    calls increment on arr[1:] to increment the next lowest digit and so on until
    it reaches the highest digit.
    When the highest digit cannot be incremented, it returns None.

    To avoid duplicate permutations, once the higher digits have been incremented
    by the recursive call, the lowest digit is set to the next lowest digit.

    A dirty hack is done to go from 2999... to 3466... when the length of the
    number is 3.

    :param arr: Number represented by an array of digits
    :param base: Arithmetic base in which this number is
    :return: The number incremented by one or None if there is no more number
    """
    if arr[0] < base - 1:
        arr[0] += 1
        return arr

    if len(arr) == 1:
        return None
    # Edge case when we hit 2999..., the next one is 3466...
    elif len(arr) == 3 and arr[1] == 6 and arr[2] == 0:
        return [3, 2, 1]

    sub = increment(arr[1:], base)
    if sub is None:
        return None
    return [sub[0]] + sub


def generate(length, base):
    """Generates a complete list of numbers of certain a length on a certain base
    that are relevant to search for the maximum multiplicative persistence.

    :param length: Number of digits of the numbers to generate
    :param base: Arithmetic base in which the numbers are
    :return: Iterator over all the relevant numbers
    """
    arr = [3] * (length-1) + [0]
    while arr:
        yield [translation[n] for n in arr]
        arr = increment(arr, base)


def multiply(n):
    """Multiply each digits of a number.

    This is optimized by splitting the int into sections of 4 digits and
    computing the multiplication on this smaller int with a cache.
    The remaining 3 or less digits are multiplied with another function that also
    has a (smaller) cache.
    """
    p = 1
    while n >= 10000:
        m = n // (10000)
        p *= _multiply_4d(n - m * 10000)
        if p == 0:
            return 0
        n = m
    return p * _multiply_3d(n)


@lru_cache(maxsize=10000)
def _multiply_4d(n):
    """Multiplication of the 4 digits of a number.
    This function is cached.
    """
    return ((n // 1000) % 10) * ((n // 100) % 10) * ((n // 10) % 10) * (n % 10)


@lru_cache(maxsize=1000)
def _multiply_3d(n):
    """Multiplication of the 4 digits of a number.
    This function is cached.
    """
    if n < 10:
        return n
    if n < 100:
        return ((n // 10) % 10) * (n % 10)
    if n < 1000:
        return ((n // 100) % 10) * ((n // 10) % 10) * (n % 10)
    return ((n // 1000) % 10) * ((n // 100) % 10) * ((n // 10) % 10) * (n % 10)


def persistence(n, steps=1):
    """Computes the multiplicative persistence of a number.
    """
    steps += 1

    p = multiply(n)
    if p < 10:
        return steps
    return persistence(p, steps)


def arr_to_str(arr):
    return "".join([str(n) for n in arr[::-1]])


def search(size):
    max_steps = 0
    n_max = 0
    # numbers_max = []
    total = 0

    for arr in generate(size, len(translation)):
        steps = persistence(reduce(lambda x, y: x * y, arr))
        if steps > max_steps:
            max_steps = steps
            n_max = 1
            # numbers_max = [arr]
        elif steps == max_steps:
            n_max += 1
            # numbers_max.append(arr)
        total += 1

        # print("{} -> {}".format(arr_to_str(arr), steps))
    # print("Total generated =", total)
    # print("Max =", max_steps)
    # print("Cache info 4digits:", _multiply_4d.cache_info())
    # print("Cache info 3digits:", _multiply_3d.cache_info())
    # print("Numbers =", [arr_to_str(arr) for arr in numbers_max])

    print("{};{};{};{}".format(
        size, max_steps, n_max, total
    ))

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: {} <number of digits>".format(sys.argv[0]))
        sys.exit(1)
    search(int(sys.argv[1]))
    # print("size;max_steps;n_max;total")
    # for n in range(2, int(sys.argv[1]) + 1):
        # search(n)
