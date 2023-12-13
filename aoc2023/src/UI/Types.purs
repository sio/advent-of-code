module UI.Types where

import Data.Maybe (Maybe)

import AOC

type State =
  { day    :: Day
  , puzzle :: Puzzle
  , result :: Maybe Solution
  , check  :: Maybe (Solution -> Boolean)
  }

data Action
  = UserInput String
  | SelectSample Int
  | SelectDay Int
  | InputTainted
