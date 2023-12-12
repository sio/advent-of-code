module Day02.Solve where

import Prelude
import Data.List (List(..), (:), foldl)
import Data.String (joinWith)
import Data.Array (zip, fromFoldable)
import Data.Generic.Rep (class Generic)
import Data.Show.Generic (genericShow)

import AOC

day02 :: Day
day02 =
  { index: 2
  , title: "Cube Conundrum"
  , solve
  , samples
  }

samples :: List Sample
samples =
  ( Sample (Numeric 8) (Empty)
    """Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green
Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue
Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red
Game 4: 1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red
Game 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green
"""
  : Nil)

solve :: Puzzle -> Solution
solve puzzle = combine (part1 puzzle) (part2 puzzle)
  where
    part1 _ = Part Nil $ Textual $ show $ (RGB 12 13 14)
    part2 _ = Part Nil Empty

data RGB = RGB Int Int Int
instance Show RGB where
  show (RGB r g b) = "RGB " <> (joinWith " " $ fromFoldable $ map show [r, g, b])
instance Semiring RGB where
  add (RGB r1 g1 b1) (RGB r2 g2 b2) = RGB (r1+r2) (g1+g2) (b1+b2)
  zero = RGB 0 0 0
  one = RGB 1 1 1
  mul (RGB r1 g1 b1) (RGB r2 g2 b2) = RGB (r1*r2) (g1*g2) (b1*b2)
instance Ring RGB where
  sub (RGB r1 g1 b1) (RGB r2 g2 b2) = RGB (r1-r2) (g1-g2) (b1-b2)
derive instance Eq RGB
instance Ord RGB where
  compare a b | a == b = EQ
  compare (RGB r1 g1 b1) (RGB r2 g2 b2) | r1 < r2 || g1 < g2 || b1 < b2 = LT
  compare _ _ = GT

data Game = Game Int (List RGB)
instance Show Game where
  show (Game index sets) = "Game " <> show index <> ": " <> showSets
    where
      showSets = joinWith "; " $ fromFoldable $ map show sets

-- | Score game based on its index and whether it's possible
score :: RGB -> Game -> Int
score limit (Game index sets) =
  if
    foldl check true sets
  then
    index
  else
    0
  where
    check :: Boolean -> RGB -> Boolean
    check false _ = false
    check _ set = set < limit
