module AOC where

import Prelude
import Data.Maybe (Maybe)

type Day =
  { index   :: Int
  , title   :: String
  , solve   :: Input -> Output
  , samples :: Array Sample
  }


newtype Input = Input String


data Output = Output Solution Debug Error
type Error = Maybe String
type Debug = Array String


data Solution = Solution Answer Answer
derive instance Eq Solution


data Answer   = Numeric Int | Textual String | Empty
derive instance Eq Answer


data Sample   = Sample Input Solution
