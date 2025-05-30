package cmd

import (
	"github.com/omni-network/omni/lib/buildinfo"
	libcmd "github.com/omni-network/omni/lib/cmd"
	"github.com/omni-network/omni/lib/errors"

	"github.com/spf13/cobra"
)

// New returns a new root cobra command that handles our command line tool.
func New() *cobra.Command {
	return libcmd.NewRootCmd(
		"omni",
		"CLI providing tools for interacting with omni",
		newOperatorCmds(),
		newDeveloperCmds(),
		newDevnetCmds(),
		buildinfo.NewVersionCmd(),
	)
}

func newOperatorCmds() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "operator",
		Short: "Operator commands",
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(
		newAVSRegisterCmd(),
		newAVSDeregisterCmd(),
		newInitCmd(),
		newCreateValCmd(),
		newCreateOperatorKeyCmd(),
		newCreateConsensusKeyCmd(),
		newUnjailCmd(),
		newDelegateCmd(),
		newEditValCmd(),
	)

	return cmd
}

func newAVSRegisterCmd() *cobra.Command {
	var cfg RegConfig

	cmd := &cobra.Command{
		Use:   "avs-register",
		Short: "Register an operator with the Omni AVS contract",
		Long: `Register command expects a Eigen-Layer yaml config file as an argument
to successfully register an operator address with the Omni AVS contract.

Note the operator must already be registered with Eigen-Layer.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			err := AVSRegister(cmd.Context(), cfg)
			if err != nil {
				return errors.Wrap(err, "registration failed")
			}

			return nil
		},
	}

	bindRegConfig(cmd, &cfg)

	return cmd
}

func newAVSDeregisterCmd() *cobra.Command {
	var cfg RegConfig

	cmd := &cobra.Command{
		Use:   "avs-deregister",
		Short: "Deregister an operator from the Omni AVS contract",
		Long: `Deregister command expects a EigenLayer yaml config file as an argument
to successfully deregister an operator from the Omni AVS contract.

Note the operator must already be registered on the Omni AVS contract.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			err := AVSDeregister(cmd.Context(), cfg)
			if err != nil {
				return errors.Wrap(err, "deregistration failed")
			}

			return nil
		},
	}

	bindRegConfig(cmd, &cfg)

	return cmd
}
