module Day01.Solve where

import Prelude
import Data.Array (toUnfoldable)
import Data.Enum (fromEnum)
import Data.Foldable (foldl)
import Data.List (List(..), (:), zip)
import Data.Maybe (Maybe(..))
import Data.String (drop, uncons, take, length, split, Pattern(..))
import Data.String.CodePoints (CodePoint, toCodePointArray, codePointFromChar)
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
  ( Sample (Numeric 142) (Numeric 142)
      """1abc2
pqr3stu8vwx
a1b2c3d4e5f
treb7uchet"""
  : Sample (Numeric 209) (Numeric 281)
      """two1nine
eightwothree
abcone2threexyz
xtwone3four
4nineeightseven2
zoneight234
7pqrstsixteen
"""
  : Sample (Numeric 13) (Numeric 12) "1four2eightseven8one3eightwogrr"
  : Nil )


solve :: Puzzle -> Solution
solve puzzle = combine (part1 puzzle) (part2 puzzle)
  where
    part1 = Part Nil <<< calibrationValue
    part2 p = Part (debug p) (calibrationParser p)

--
-- Part 1: A simple fold
--

data State = State (Maybe Int) (Maybe Int) Int

calibrationValue :: String -> Answer
calibrationValue = Numeric <<< unpack <<< iterate
  where
    initial :: State
    initial = State Nothing Nothing 0

    unpack :: State -> Int
    unpack (State Nothing _ x) = x
    unpack (State _ Nothing x) = x
    unpack (State (Just first) (Just last) sum) = sum + first * 10 + last

    iterate :: String -> State
    iterate = foldl parse initial <<< toCodePointArray

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

--
-- Part 2: String parsing
--
-- There exist proper approaches to parsing in Purescript,
-- but I didn't learn those yet
--
-- Let's try to push through with what I know so far...

calibrationParser :: String -> Answer
calibrationParser input = Numeric $ total $ foldl worker init tokens
  where
    tokens = parser (Cursor input Nil)

    init = {first: 0, last: Nothing, sum: 0}

    total {first, last: Just last', sum} = sum + first * 10 + last'
    total {last: Nothing, sum} = sum

    worker state Newline = init {sum = total state}
    worker {last: Nothing, sum} (Digit d) = {first: d, last: Just d, sum}
    worker {last, sum} (Digit d) = {first: d, last, sum}
    worker state _ = state {sum = state.sum + errMarker}

    errMarker = 5003 -- large prime number easy to sum in your head

debug :: String -> Log
debug input = map Info $ map (\(Tuple line token) -> line <> " -> " <> show token) $ zip lines tokens
  where
    lines = toUnfoldable $ split (Pattern "\n") input

    tokens :: List (List Token)
    tokens = group $ parser (Cursor input Nil)

    group :: List Token -> List (List Token)
    group = foldl go init
      where
        init :: List (List Token)
        init = Nil

        go :: List (List Token) -> Token -> List (List Token)
        go Nil Newline = Nil
        go acc Newline = (Nil : acc)
        go Nil token = ((token : Nil) : Nil)
        go (x : xs) token = ((token : x) : xs)

data Token = Digit Int | Newline | ParsingError
instance Show Token where
  show Newline = "LF"
  show ParsingError = "ERR"
  show (Digit x) = show x

data Cursor = Cursor String (List Token)

-- Tokens are produced in reverse order
parser :: Cursor -> List Token
parser (Cursor "" tokens) = tokens
parser (Cursor text tokens) | hasDigit text =
  parser $ Cursor (drop 1 text) (getDigit text : tokens)
parser (Cursor text tokens) | hasCursive text =
  parser $ Cursor (drop 1 text) (getCursive text : tokens)
parser (Cursor text tokens) | take 1 text == "\n" =
  parser $ Cursor (drop 1 text) (Newline : tokens)
parser (Cursor text tokens) =
  parser $ Cursor (drop 1 text) tokens

hasDigit :: String -> Boolean
hasDigit s = case parseDigit s of
  Nothing -> false
  _ -> true

getDigit :: String -> Token
getDigit s = case parseDigit s of
  Nothing -> ParsingError
  Just d -> Digit d

parseDigit :: String -> Maybe Int
parseDigit s = case uncons s of
  Nothing -> Nothing
  Just {head} ->
    if
      delta < 10 && delta >= 0
    then
      Just delta
    else
      Nothing
    where
      delta = (fromEnum head) - (fromEnum $ codePointFromChar '0')

hasCursive :: String -> Boolean
hasCursive s = case getCursive s of
  ParsingError -> false
  _ -> true

cursiveDigits :: List String
cursiveDigits = toUnfoldable
  ["zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"]

getCursive :: String -> Token
getCursive text = value $ foldl worker init cursiveDigits
  where
    value :: Tuple Int Token -> Token
    value (Tuple _ token) = token

    init :: Tuple Int Token
    init = Tuple 0 ParsingError

    worker :: (Tuple Int Token) -> String -> (Tuple Int Token)
    worker (Tuple index ParsingError) cursive =
      if
        (take len text) == cursive
      then
        Tuple 0 (Digit index)
      else
        Tuple (index+1) ParsingError
      where
        len = length cursive
    worker digit _ = digit
