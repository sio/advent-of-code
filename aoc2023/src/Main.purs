module Main where

import Prelude
import Effect (Effect)
import Halogen.Aff as HA
import Halogen.VDom.Driver (runUI)

import AOC (Day)
import UI.Logic (aoc)

import Day01.Solve (day01)
import Day02.Solve (day02)

days :: Array Day
days =
  [ day01
  , day02
  ]

main :: Effect Unit
main = HA.runHalogenAff do
  body <- HA.awaitBody
  runUI (aoc days) unit body
