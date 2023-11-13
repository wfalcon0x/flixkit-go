package v1_0_0

import (
	"context"
	"fmt"

	"github.com/onflow/cadence/runtime/ast"
	"github.com/onflow/cadence/runtime/parser"
	"github.com/onflow/flixkit-go"
	"github.com/onflow/flixkit-go/generator"
	"github.com/onflow/flow-cli/flowkit"
	"github.com/onflow/flow-cli/flowkit/config"
	"github.com/onflow/flow-cli/flowkit/gateway"
	"github.com/onflow/flow-cli/flowkit/output"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/spf13/afero"
)

type GeneratorV1_0_0 struct {
	deployedContracts []flixkit.Contracts
	coreContracts     flixkit.Contracts
	testnetClient     *flowkit.Flowkit
	mainnetClient     *flowkit.Flowkit
}

// stubb to pass in parameters
func NewGenerator(deployedContracts []flixkit.Contracts, coreContracts flixkit.Contracts, logger output.Logger) (flixkit.Generator, error) {
	loader := afero.Afero{Fs: afero.NewOsFs()}

	gwt, err := gateway.NewGrpcGateway(config.TestnetNetwork)
	if err != nil {
		return nil, fmt.Errorf("could not create grpc gateway for testnet %w", err)
	}

	gwm, err := gateway.NewGrpcGateway(config.MainnetNetwork)
	if err != nil {
		return nil, fmt.Errorf("could not create grpc gateway for mainnet %w", err)
	}

	state, err := flowkit.Init(loader, crypto.ECDSA_P256, crypto.SHA3_256)
	if err != nil {
		return nil, fmt.Errorf("could not initialize flowkit state %w", err)
	}
	testnetClient := flowkit.NewFlowkit(state, config.TestnetNetwork, gwt, logger)
	mainnetClient := flowkit.NewFlowkit(state, config.MainnetNetwork, gwm, logger)

	if coreContracts == nil {
		coreContracts = generator.GetDefaultCoreContracts()
	}

	return &GeneratorV1_0_0{
		deployedContracts: deployedContracts,
		coreContracts:     coreContracts,
		testnetClient:     testnetClient,
		mainnetClient:     mainnetClient,
	}, nil
}

func (g GeneratorV1_0_0) Generate(ctx context.Context, code string, preFill *flixkit.FlowInteractionTemplate) (*flixkit.FlowInteractionTemplate, error) {
	template := &flixkit.FlowInteractionTemplate{}
	if preFill != nil {
		template = preFill
	}

	// make sure imports use new import syntax "string import"
	normalizedCode := generator.NormalizeImports(code)

	codeBytes := []byte(normalizedCode)
	program, err := parser.ParseProgram(nil, codeBytes, parser.Config{})
	if err != nil {
		return nil, err
	}

	err = generator.ProcessParameters(program, template)
	if err != nil {
		return nil, err
	}

	// save v1.0.0 cadence code to template, with placeholder imports
	template.Data.Cadence = generator.UnNormalizeImports(normalizedCode)
	err = g.processDependencies(ctx, program, template)
	if err != nil {
		return nil, err
	}

	// ignore interface type for now
	template.FType = "InteractionTemplate"
	template.FVersion = "1.0.0"
	template.Data.Type = generator.DetermineCadenceType(program)
	id, err := flixkit.GenerateFlixID(template)
	if err != nil {
		return nil, err
	}
	template.ID = id

	return template, nil
}

func (g GeneratorV1_0_0) processDependencies(ctx context.Context, program *ast.Program, template *flixkit.FlowInteractionTemplate) error {
	imports := program.ImportDeclarations()

	if len(imports) == 0 {
		return nil
	}

	// fill in dependence information
	deps := make(flixkit.Dependencies, len(imports))
	for _, imp := range imports {
		contractName, err := generator.ExtractContractName(imp.String())
		if err != nil {
			return err
		}
		dep, err := g.generateDependenceInfo(ctx, contractName)
		if err != nil {
			return err
		}
		for contractName, contract := range dep {
			deps[contractName] = contract
		}
		template.Data.Dependencies = deps
	}

	return nil
}

func (g *GeneratorV1_0_0) generateDependenceInfo(ctx context.Context, contractName string) (map[string]flixkit.Contracts, error) {
	var placeholder string
	var info flixkit.Networks
	// new import syntax detected, convert to old import syntax, limitation of 1.0.0
	placeholder = "0x" + contractName
	info = generator.GetContractInformation(contractName, g.deployedContracts, g.coreContracts)

	for name, network := range info {
		var flowkit *flowkit.Flowkit
		if name == config.MainnetNetwork.Name && g.mainnetClient != nil {
			flowkit = g.mainnetClient
		} else if name == config.TestnetNetwork.Name && g.testnetClient != nil {
			flowkit = g.testnetClient
		}
		if network.Pin == "" && flowkit != nil {
			hash, height, _ := generator.GeneratePinDebthFirst(ctx, *flowkit, network.Address, network.Contract)
			network.Pin = hash
			network.PinBlockHeight = height
		}
		info[name] = network
	}

	if info == nil {
		return nil, fmt.Errorf("contract %s not found", contractName)
	}

	return map[string]flixkit.Contracts{
		placeholder: {
			contractName: info,
		},
	}, nil
}
