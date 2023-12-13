module UI.Types where

import Data.Maybe (Maybe)

import AOC

type State =
  { day    :: Day
  , puzzle :: Puzzle
  , result :: Solution
  , check  :: Maybe (Solution -> Boolean)
  }

data Action = UserInput String | SelectSample Int | SelectDay Int
