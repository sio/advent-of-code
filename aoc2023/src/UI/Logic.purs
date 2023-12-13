module UI.Logic where

import Prelude
import Data.Maybe (Maybe(..), fromMaybe)
import Data.List (List(..), (!!))
import Data.Array (last, find)
import Halogen as H

import AOC
import UI.Types
import UI.Layout (render)

aoc :: forall query input output m. Array Day -> H.Component query input output m
aoc days =
  H.mkComponent
    { initialState: initialize days
    , render: render days
    , eval: H.mkEval H.defaultEval { handleAction = handle days }
    }

initialize :: forall input. Array Day -> input -> State
initialize days _ = set day 0
  where
    day = fromMaybe zeroDay $ last days
    zeroDay =
      { index: 0
      , title: "Empty day"
      , solve: \_ -> Solution Nil Empty Empty
      , samples: Nil
      }

set :: Day -> Int -> State
set day sampleIndex =
  { day
  , puzzle
  , result: day.solve puzzle
  , check:  map match sample
  } where
      sample = day.samples !! sampleIndex
      puzzle = case sample of
                Nothing -> ""
                Just (Sample _ _ i) -> i

handle :: forall output m. Array Day -> Action -> H.HalogenM State Action () output m Unit
handle _ (UserInput s) =
  H.modify_ \state -> state { puzzle = s, result = state.day.solve s, check = Nothing }
handle _ (SelectSample n) =
  H.modify_ \state -> (set state.day n)
handle days (SelectDay n) =
  H.modify_ \state ->
    let
      day = find (\d -> d.index == n) days
    in case day of
      Nothing -> state
      Just day -> set day 0
