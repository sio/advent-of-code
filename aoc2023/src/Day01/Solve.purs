module Day01.Solve (day01) where

import Data.Maybe (Maybe(..))

import AOC

day01 :: Day
day01 =
  { index: 0
  , title: "Sample AoC solution"
  , solve
  , samples
  }

samples :: Array Sample
samples = []

solve :: Input -> Output
solve _ = Output (Solution Empty Empty) [] Nothing
