module Day01.Solve where

import Prelude
import Data.Maybe (Maybe(..), fromMaybe)
import Data.List (List(..), (:))
import Data.String (uncons)
import Data.String.CodePoints (codePointFromChar)
import Data.Enum (fromEnum)
import Data.Tuple (Tuple(..))


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
  : Nil )


int :: Maybe Int -> Int
int x = fromMaybe 0 x

solve :: Puzzle -> Solution
solve puzzle = combine (part1 puzzle) (part2 puzzle)
  where
    part1 p = Part Nil $ Numeric $ textValue p
    part2 _ = Part Nil Empty

textValue = go Nothing Nothing 0 where
  go :: Maybe Int -> Maybe Int -> Int -> String -> Int
  go (Just first) (Just last) total "" = total + first * 10 + last
  go Nothing Nothing total "" = total
  go Nothing _ total text =
    let
      ParseResult first _ tail = headDigit text
    in
      go first Nothing total tail
  go first Nothing total text =
    let
      ParseResult last _ tail = headDigit text
    in
      go first last total tail
  go first last total text =
    case headDigit text of
      ParseResult Nothing _  tail -> go first last total tail
      ParseResult next false tail -> go first next total tail
      ParseResult next true  tail -> go next Nothing (total + (int first) * 10 + int last) tail

data ParseResult = ParseResult (Maybe Int) Boolean String
headDigit :: String -> ParseResult
headDigit = go false where
  go :: Boolean -> String -> ParseResult
  go newline line =
    case (uncons line) of
      Nothing -> ParseResult Nothing newline ""
      Just {head, tail} ->
        let
          digit = (fromEnum head) - (fromEnum (codePointFromChar '0'))
          br = (head == codePointFromChar '\n' || newline)
        in
          if
            digit < 10 && digit >= 0
          then
            ParseResult (Just digit) br tail
          else
            go br tail
