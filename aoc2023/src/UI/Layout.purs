module UI.Layout (render) where

import Prelude
import Data.Array (fromFoldable)
import Data.List (List(..))
import Data.Maybe (Maybe(..))
import Data.Traversable (scanl)
import Halogen as H
import Halogen.HTML as HH
import Halogen.HTML.Events as HE
import Halogen.HTML.Properties as HP

import AOC
import UI.Types

render :: forall m. Array Day -> State -> H.ComponentHTML Action () m
render days state =
  HH.main_
    [ title
    , navigation days
    , header state.day
    , samples state.day
    , HH.textarea
        [ HE.onValueChange UserInput
        , HP.value state.puzzle
        , classes ["puzzle"]
        ]
    , solution state
    , footer
    ]

title :: HH.HTML _ _
title = HH.h1_ [HH.text "Advent of Code in Purescript" ]

navigation :: Array Day -> HH.HTML _ _
navigation days = HH.nav_ $ map button days
  where
    button d = HH.button [HE.onClick (\_ -> SelectDay d.index)] [HH.text $ "Day " <> show d.index]

header :: Day -> HH.HTML _ _
header day = HH.header_
  [ HH.h2_ [HH.text $ "Day " <> show day.index <> ": " <> day.title]
  , HH.a [HP.href dayUrl] [HH.text "Puzzle description"]
  , HH.text ", "
  , HH.a [HP.href puzzleUrl] [HH.text "personalized input file"]
  ]
  where
    prefix = "https://adventofcode.com/2023/day/"
    dayUrl = prefix <> show day.index
    puzzleUrl = dayUrl <> "/input"

footer :: HH.HTML _ _
footer = HH.footer_ [HH.a [HP.href github] [HH.text github]]
  where
    github = "https://github.com/sio/advent-of-code"

samples :: Day -> HH.HTML _ _
samples day =
  HH.nav_
    $ map button
    $ fromFoldable
    $ scanl (\x _ -> x + 1) 0 day.samples
  where
    button i =
      HH.button [HE.onClick (\_ -> SelectSample i)] [HH.text $ "Sample " <> show i]

solution :: State -> HH.HTML _ _
solution state = result state.result
  where
    result (Solution log part1 part2) = HH.div_
      [ answerContainer 1 part1
      , answerContainer 2 part2
      , renderCheck state.check
      , renderLog log
      ]

    answerContainer index answer = HH.div [classes ["answer"]]
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
      , HH.ul [classes ["log"]] logLines
      ]
      where
        logLines = fromFoldable $ map renderLine l
        renderLine line = HH.li_ [HH.text $ show line]

    renderCheck Nothing = HH.text ""
    renderCheck (Just check) = HH.div
      [classes ["sample-match-" <> show c]]
      [HH.text $ if c then "OK" else "FAIL"]
      where c = check state.result

classes :: forall r i. Array String -> HP.IProp (class :: String | r) i
classes = HP.classes <<< map HH.ClassName
