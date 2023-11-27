module AOC where

import Prelude
import Data.List (List(..))

type Day =
  { index   :: Int
  , title   :: String
  , solve   :: Input -> Output
  , samples :: List Sample
  }

newtype Input = Input String
data Output = Output Solution Debug Status

data Status = Ok | Error String
type Debug = List String

data Solution = Solution Answer Answer
derive instance Eq Solution

data Answer   = Numeric Int | Textual String | Empty
derive instance Eq Answer

data Sample   = Sample Input Solution
