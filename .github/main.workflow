workflow "go fmt all" {
  on = "pull_request"
  resolves = ["go fmt"]
}

action "go fmt" {
  uses    = "sjkaliski/go-github-actions/fmt@v0.4.0"
  secrets = ["GITHUB_TOKEN"]
}

workflow "go lint all" {
  on = "pull_request"
  resolves = ["go lint"]
}

action "go lint" {
  uses    = "sjkaliski/go-github-actions/lint@v0.4.0"
  secrets = ["GITHUB_TOKEN"]
}

workflow "go test all" {
  on = "push"
  resolves = ["go test"]
}

action "go test" {
  uses = "docker://docker.io/golang:1.12"
  runs = "go"
  args = "test ./..."
}
