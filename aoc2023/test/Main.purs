module Test.Main where

import Prelude

import Data.Array (toUnfoldable)
import Data.List (List, (:))
import Data.Traversable (foldl)
import Effect (Effect)
import Effect.Class.Console (log)
import Test.Assert (assert)

import Main (days)
import AOC (Day, LogEntry(..), Sample(..), match, Log)

main :: Effect Unit
main = do
  log stdout
  assert testsOK
    where
      stdout :: String
      stdout = foldl (\prev cur -> prev <> "\n" <> show cur) "" tests

      testsOK :: Boolean
      testsOK = foldl testSuccess true tests

      testSuccess false _ = false
      testSuccess _ (Error _) = false
      testSuccess t _ = t

      tests :: Log
      tests = toUnfoldable $ map checkSamples days

checkSamples :: Day -> LogEntry
checkSamples day = joinLog $ ((Info $ "Day " <> show day.index) : (map check day.samples))
  where
    check sample@(Sample _ _ puzzle) =
      if match sample (day.solve puzzle)
      then Info "OK"
      else Error "FAIL"

joinLog :: List LogEntry -> LogEntry
joinLog = foldl go (Info "")
  where
    go (Info "") cur = cur
    go (Info prev) (Info cur) = Info (prev <> "; " <> cur)
    go (Info prev) (Error cur) = Error (prev <> "; " <> cur)
    go (Error prev) (Info cur) = Error (prev <> "; " <> cur)
    go (Error prev) (Error cur) = Error (prev <> "; " <> cur)
