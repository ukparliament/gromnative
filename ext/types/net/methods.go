package net

type Methods interface {
  Get (*GetInput) (*GetOutput, error)
}