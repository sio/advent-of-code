module Day02.Solve where

import Prelude
import Control.Alt ((<|>))
import Data.Array (zip, fromFoldable)
import Data.Either (fromRight)
import Data.Generic.Rep (class Generic)
import Data.List (List(..), (:), foldl)
import Data.Show.Generic (genericShow)
import Data.String (joinWith)
import Parsing (Parser, runParser)
import Parsing.Combinators (many, many1, optional)
import Parsing.String (string, char)
import Parsing.String.Basic (intDecimal, space)

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
  ( Sample (Numeric 8) (Numeric 2286)
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
    games = fromRight Nil $ runParser puzzle (many game)
    part1 puzzle =
      Part log
        $ Numeric
        $ foldl add zero
        $ map (score ceiling)
        $ games
      where
        ceiling = RGB 12 13 14
        log = map debug games
        debug game = Info $ (show game) <> " -> " <> (show $ score ceiling game)
    part2 _ =
      Part Nil
        $ Numeric
        $ foldl add zero
        $ map power
        $ map (\(Game _ sets) -> fewest sets) games

-- | Part 2: Fewest cubes required to play the game
fewest :: List RGB -> RGB
fewest = foldl ceil mempty
  where
    ceil :: RGB -> RGB -> RGB
    ceil (RGB r1 g1 b1) (RGB r2 g2 b2) = RGB (max r1 r2) (max g1 g2) (max b1 b2)

-- | Part 2: Scoring function
power :: RGB -> Int
power (RGB r g b) = r * g * b

-- | Single set of RGB cubes
data RGB = RGB Int Int Int

instance Show RGB where
  show (RGB r g b) = "RGB " <> (joinWith " " $ fromFoldable $ map show [r, g, b])

derive instance Eq RGB

fits :: RGB -> RGB -> Boolean
fits (RGB limR limG limB) (RGB r g b) = r <= limR && g <= limG && b <= limB

instance Monoid RGB where
  mempty = RGB 0 0 0

instance Semigroup RGB where
  append (RGB r1 g1 b1) (RGB r2 g2 b2) = RGB (r1+r2) (g1+g2) (b1+b2)

-- | Full cubes game state
data Game = Game Int (List RGB)

instance Show Game where
  show (Game index sets) = "Game " <> show index <> ": " <> showSets
    where
      showSets = joinWith "; " $ fromFoldable $ map show sets

-- | Add a set of cubes to the game
draw :: Game -> RGB -> Game
draw (Game index sets) cubes = Game index (cubes : sets)

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
    check _ set = fits limit set

-- | Parse game input
game :: Parser String Game
game = do
  _ <- string "Game"
  _ <- many1 space
  index <- intDecimal
  _ <- string ": "
  draws <- many rgb
  pure $ Game index draws

rgb :: Parser String RGB
rgb = do
  cubes <- many1 cube
  _ <- optional $ char ';'
  _ <- many space
  pure $ foldl append mempty cubes

cube :: Parser String RGB
cube = do
  count <- intDecimal
  _ <- many1 space
  which <- color
  _ <- optional $ char ','
  _ <- many space
  pure $ which count

data Color = Red | Green | Blue
color :: Parser String (Int -> RGB)
color = do
  name <- (string "red") <|> (string "green") <|> (string "blue")
  _ <- many $ string ","
  _ <- many space
  let tag =
        if name == "red" then
          Red else
        if name == "green" then
          Green
        else
          Blue
  pure $ paint tag

-- | Paint a number of cubes with given color
paint :: Color -> Int -> RGB
paint Red x   = RGB x 0 0
paint Green x = RGB 0 x 0
paint Blue x  = RGB 0 0 x
