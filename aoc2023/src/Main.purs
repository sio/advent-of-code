module Main where

import Prelude
import Data.List (List(..), (!!))
import Data.Maybe (Maybe(..))
import Data.Traversable (scanl)
import Data.Array (fromFoldable, find)

import Effect (Effect)
import Halogen as H
import Halogen.Aff as HA
import Halogen.HTML as HH
import Halogen.HTML.Events as HE
import Halogen.HTML.Properties as HP
import Halogen.VDom.Driver (runUI)

import AOC

import Day01.Solve (day01)
import Day02.Solve (day02)

days :: Array Day
days =
  [ day01
  , day02
  ]

type State =
  { day    :: Day
  , puzzle :: Puzzle
  , result :: Solution
  , check  :: Maybe (Solution -> Boolean)
  }

initialState :: forall input. input -> State
initialState _ = setState day01 0

setState :: Day -> Int -> State
setState day sampleIndex =
  { day
  , puzzle
  , result: day.solve puzzle
  , check:  map match sample
  } where
      sample = day.samples !! sampleIndex
      puzzle = case sample of
                Nothing -> ""
                Just (Sample _ _ i) -> i

data Action = UserInput String | SelectSample Int | SelectDay Int

handleAction :: forall output m. Action -> H.HalogenM State Action () output m Unit
handleAction (UserInput s) =
  H.modify_ \state -> state { puzzle = s, result = state.day.solve s, check = Nothing }
handleAction (SelectSample n) =
  H.modify_ \state -> (setState state.day n)
handleAction (SelectDay n) =
  H.modify_ \state ->
    let
      day = find (\d -> d.index == n) days
    in case day of
      Nothing -> state
      Just day -> setState day 0

render :: forall m. State -> H.ComponentHTML Action () m
render state =
  HH.main_
    [ renderHeader state.day
    , renderSamples state.day
    , HH.textarea
        [ HE.onValueChange UserInput
        , HP.value state.puzzle
        , HP.classes [ HH.ClassName "puzzle" ]
        ]
    , renderSolution state.result
    , renderFooter
    ]
  where
    renderHeader d = HH.header_
      [ HH.h1_ [HH.text "Advent of Code in Purescript" ]
      , renderDays
      , HH.h2_ [HH.text $ "Day " <> show d.index <> ": " <> d.title]
      , HH.a [HP.href dayUrl] [HH.text "Puzzle description"]
      , HH.text ", "
      , HH.a [HP.href puzzleUrl] [HH.text "personalized input file"]
      ]
      where
        prefix = "https://adventofcode.com/2023/day/"
        dayUrl = prefix <> show d.index
        puzzleUrl = dayUrl <> "/input"

    renderFooter = HH.footer_
      [HH.a [HP.href github] [HH.text github]]
      where
        github = "https://github.com/sio/advent-of-code"

    renderSamples d = HH.div_ $
      map renderSample $ fromFoldable $ scanl (\x _ -> x + 1) 0 d.samples
    renderSample i =
      HH.button [HE.onClick (\_ -> SelectSample (i-1))] [HH.text $ "Sample " <> show i]

    renderDays = HH.div_ $ map dayButton days
    dayButton d =
      HH.button [HE.onClick (\_ -> SelectDay d.index)] [HH.text $ "Day " <> show d.index]

    renderSolution (Solution log part1 part2) = HH.div_
      [ answerContainer 1 part1
      , answerContainer 2 part2
      , renderCheck state.check
      , renderLog log
      ]

    answerContainer index answer = HH.div [HP.classes [HH.ClassName "answer"]]
      [ HH.span_ [HH.text $ "Part " <> show index <> ": " ]
      , HH.span_ [renderAnswer answer]
      ]

    renderAnswer Empty = HH.text "Not solved"
    renderAnswer (Numeric n) = HH.text $ show n
    renderAnswer (Textual t) = HH.pre_ [HH.text t]

    renderLog Nil = HH.text ""
    renderLog l =
      HH.details_
      [ HH.summary_ [HH.text "Debug log"]
      , HH.ul [HP.classes [HH.ClassName "log"]] logLines
      ]
      where
        logLines = fromFoldable $ map renderLine l
        renderLine line = HH.li_ [HH.text $ show line]

    renderCheck Nothing = HH.text ""
    renderCheck (Just check) = HH.div
      [HP.classes [HH.ClassName $ "sample-match-" <> show c]]
      [HH.text $ if c then "OK" else "FAIL"]
      where c = check state.result

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
