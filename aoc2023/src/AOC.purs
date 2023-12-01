module AOC where

import Prelude
import Data.List (List(..), (:))

type Day =
  { index   :: Int
  , title   :: String
  , solve   :: Puzzle -> Solution
  , samples :: List Sample
  }

type Puzzle = String

data Solution = Solution Log Answer Answer
instance Eq Solution where
  eq (Solution _ a1 a2) (Solution _ b1 b2) = (a1 == b1) && (a2 == b2)

combine :: Part -> Part -> Solution
combine (Part log1 ans1) (Part log2 ans2) = Solution (log1 <> log2) ans1 ans2

data Part = Part Log Answer
instance Eq Part where
  eq (Part log1 ans1) (Part log2 ans2) =
    if
      hasError log1 || hasError log2
    then
      false
    else
      ans1 == ans2

data Answer   = Numeric Int | Textual String | Empty
derive instance Eq Answer

data Sample   = Sample Answer Answer Puzzle
match :: Sample -> Solution -> Boolean
match (Sample want1 want2 _) (Solution _ got1 got2) = (want1 == got1) && (want2 == got2)

type Log = List LogEntry
data LogEntry = Info String | Error String
instance Show LogEntry where
  show (Info s) = "INFO: " <> s
  show (Error s) = "ERROR: " <> s

hasError :: Log -> Boolean
hasError Nil = false
hasError (Error _ : _) = true
hasError (_ : xs) = hasError xs
