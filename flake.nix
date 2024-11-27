{
  description = "Flake for setting up the environment for the Mediator package";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
    pre-commit-hooks.url = "github:cachix/pre-commit-hooks.nix";
    devshell.url = "github:numtide/devshell";
  };

  outputs = inputs @ {flake-parts, ...}:
    flake-parts.lib.mkFlake {inherit inputs;} {
      systems = ["x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin"];
      imports = [
        inputs.pre-commit-hooks.flakeModule
        inputs.devshell.flakeModule
      ];

      perSystem = {
        config,
        self',
        inputs',
        pkgs,
        system,
        ...
      }: {
        formatter = pkgs.alejandra;

        pre-commit = {
          settings = {
            hooks = {
              alejandra = {
                enable = true;
                excludes = ["vendor"];
              };
              golangci-lint = {
                enable = true;
              };
              gotest = {
                enable = false;
                excludes = ["vendor" "mocks"];
              };
              govet = {
                enable = true;
                excludes = ["vendor" "mocks"];
              };
            };
          };
        };

        devshells.default = {
          packages = with pkgs; [
            go
            golangci-lint
            gnumake
          ];

          devshell.startup = {
            preCommitHooks.text = config.pre-commit.installationScript;
            init.text = ''
              ${pkgs.go}/bin/go mod tidy
            '';
          };

          commands = [
            {
              help = "Build the code";
              name = "build";
              category = "go";
              command = "${pkgs.go}/bin/go build ./...";
            }
            {
              help = "Check the go source code";
              name = "vet";
              category = "go";
              command = "${pkgs.go}/bin/go vet ./...";
            }
            {
              help = "Run the tests";
              name = "test";
              category = "go";
              command = "${pkgs.go}/bin/go test -shuffle=on ./...";
            }
            {
              help = "Calculate the test coverage";
              name = "cover";
              category = "go";
              command = ''
                ${pkgs.go}/bin/go test -coverprofile=coverage.out -covermode=atomic .
                ${pkgs.go}/bin/go tool cover -html=coverage.out -o coverage.html
              '';
            }
            {
              help = "(Re)generate files";
              name = "gen";
              category = "go";
              command = "${pkgs.go}/bin/go generate ./...";
            }
            {
              help = "Lint using golangci-lint";
              name = "lint";
              category = "go";
              command = "${pkgs.golangci-lint}/bin/golangci-lint run -c .golangci.yml";
            }
          ];
        };
      };
    };
}
