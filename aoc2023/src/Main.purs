module Main where

import Prelude
import Data.List ((!!))
import Data.Maybe (Maybe(..))

import Effect (Effect)
import Halogen as H
import Halogen.Aff as HA
import Halogen.HTML as HH
import Halogen.HTML.Events as HE
import Halogen.HTML.Properties as HP
import Halogen.VDom.Driver (runUI)

import AOC

import Day01.Solve (day01)

days :: Array Day
days =
  [ day01
  ]

type State =
  { day    :: Day
  , puzzle :: Puzzle
  , result :: Solution
  , check  :: Solution -> Boolean
  }

initialState :: forall input. input -> State
initialState _ = setState day01 0

setState :: Day -> Int -> State
setState day sampleIndex =
  { day
  , puzzle
  , result: day.solve puzzle
  , check:  case sample of
              Nothing -> const true
              Just s  -> match s
  } where
      sample = day.samples !! sampleIndex
      puzzle = case sample of
                Nothing -> ""
                Just (Sample _ _ i) -> i

data Action = UserInput String

handleAction :: forall output m. Action -> H.HalogenM State Action () output m Unit
handleAction (UserInput s) =
  H.modify_ \state -> state { puzzle = s, result = state.day.solve s, check = const true }

render :: forall m. State -> H.ComponentHTML Action () m
render state =
  HH.main_
    [ renderHeader state.day
    , HH.textarea
        [ HE.onValueChange UserInput
        , HP.value state.puzzle
        , HP.classes [ HH.ClassName "puzzle" ]
        ]
    , renderSolution state.result
    , renderCheck $ state.check state.result
    ]
  where
    renderHeader d = HH.header_
      [ HH.h1_ [HH.text "Advent of Code in Purescript" ]
      , HH.h2_ [HH.text $ "Day " <> show d.index <> ": " <> d.title]
      ]
    renderSolution (Solution log part1 part2) = HH.div_
      [ renderAnswer part1
      , renderAnswer part2
      , renderLog log
      ]
    renderAnswer Empty = HH.div_ [HH.text "Not solved"]
    renderAnswer (Numeric n) = HH.div
      [ HP.classes
        [ HH.ClassName "numeric"
        , HH.ClassName "answer"
        ]
      ]
      [ HH.text $ show n ]
    renderAnswer (Textual t) = HH.div
      [ HP.classes
        [ HH.ClassName "textual"
        , HH.ClassName "answer"
        ]
      ]
      [ HH.code_ [HH.text t] ]
    renderLog l = HH.code_ [HH.text $ show l]
    renderCheck c = HH.div
      [HP.classes [HH.ClassName $ "sample-match-" <> show c]]
      [HH.text $ if c then "OK" else "FAIL"]

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
