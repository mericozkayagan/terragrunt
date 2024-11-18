terraform {
  source = "."

  before_hook "tflint" {
    commands = ["apply", "plan"]
    execute = ["tflint"]
  }
}

