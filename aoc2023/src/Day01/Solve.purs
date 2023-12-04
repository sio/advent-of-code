module Day01.Solve where

import Prelude
import Data.Foldable (foldl)
import Data.Maybe (Maybe(..))
import Data.List (List(..), (:))
import Data.String.CodePoints (CodePoint, toCodePointArray, codePointFromChar)
import Data.Enum (fromEnum)


import AOC

day01 :: Day
day01 =
  { index: 1
  , title: "Trebuchet?!"
  , solve
  , samples
  }

samples :: List Sample
samples =
  ( Sample (Numeric 142) Empty
      """1abc2
pqr3stu8vwx
a1b2c3d4e5f
treb7uchet"""
  : Sample Empty (Numeric 281)
      """two1nine
eightwothree
abcone2threexyz
xtwone3four
4nineeightseven2
zoneight234
7pqrstsixteen
"""
  : Nil )


solve :: Puzzle -> Solution
solve puzzle = combine (part1 puzzle) (part2 puzzle)
  where
    part1 p = Part Nil $ calibrationValue p
    part2 _ = Part Nil Empty

data State = State (Maybe Int) (Maybe Int) Int

calibrationValue :: String -> Answer
calibrationValue puzzle = Numeric $ unpack $ last puzzle
  where
    initial :: State
    initial = State Nothing Nothing 0

    unpack :: State -> Int
    unpack (State Nothing _ x) = x
    unpack (State _ Nothing x) = x
    unpack (State (Just first) (Just last) sum) = sum + first * 10 + last

    last :: String -> State
    last s = foldl parse initial (toCodePointArray s)

    digit :: CodePoint -> Maybe Int
    digit c = if delta < 10 && delta >= 0 then Just delta else Nothing
      where
        delta = (fromEnum c) - (fromEnum $ codePointFromChar '0')

    parse :: State -> CodePoint -> State
    parse state c | c == codePointFromChar '\n' = State Nothing Nothing (unpack state)
    parse (State Nothing _ sum) c = State (digit c) (digit c) sum
    parse (State first last sum) c =
      let
        d = digit c
      in case d of
        Nothing -> State first last sum
        _ -> State first d sum

