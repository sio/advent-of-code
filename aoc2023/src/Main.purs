module Main where

import Prelude
import Data.List (head)
import Data.Maybe (Maybe(..), fromMaybe)

import Effect (Effect)
import Halogen as H
import Halogen.Aff as HA
import Halogen.HTML as HH
import Halogen.HTML.Events as HE
import Halogen.HTML.Properties as HP
import Halogen.VDom.Driver (runUI)

import AOC

import Day01.Solve

days =
  [ day01
  ]

type State =
  { day :: Day
  , puzzle :: Input
  , result :: Output
  , expect :: Maybe Solution
  }

initialState :: forall input. input -> State
initialState _ =
  { day: day
  , puzzle: puzzle
  , result: day.solve puzzle
  , expect: case sample of
              Nothing -> Nothing
              Just (Sample _ solution) -> Just solution
  } where
      day = day01
      sample = head day.samples
      puzzle = case sample of
                Nothing -> ""
                Just (Sample i _) -> i

data Action = UserInput String

render :: forall m. State -> H.ComponentHTML Action () m
render state =
  HH.main_
    [ HH.textarea
        [ HE.onValueChange UserInput
        , HP.value state.puzzle
        ]
    , HH.div_ [HH.text state.puzzle]
    ]

handleAction :: forall output m. Action -> H.HalogenM State Action () output m Unit
handleAction (UserInput s) =
  H.modify_ \state -> state { puzzle = s, expect = Nothing }

component :: forall query input output m. H.Component query input output m
component =
  H.mkComponent
    { initialState
    , render
    , eval: H.mkEval H.defaultEval { handleAction = handleAction }
    }

main :: Effect Unit
main = HA.runHalogenAff do
  body <- HA.awaitBody
  runUI component unit body
