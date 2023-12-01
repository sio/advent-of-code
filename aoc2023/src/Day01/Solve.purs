module Day01.Solve (day01) where

import Data.List (List(..))

import AOC

day01 :: Day
day01 =
  { index: 1
  , title: "Sample AoC solution"
  , solve
  , samples
  }

samples :: List Sample
samples = Nil

solve :: Puzzle -> Solution
solve _ = Solution Nil Empty Empty
