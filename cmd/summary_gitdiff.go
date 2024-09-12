package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/vx416/smart-git/ai"
	"github.com/vx416/smart-git/git"
)

var (
	// SummaryGitDiff is the command to summarize the git diff.
	SummaryGitDiff = &cobra.Command{
		Use:     "summary-gitdiff",
		Aliases: []string{"sgd"},
		Short:   "Summarize the git diff",
		Example: ` smart-git summary-gitdiff --input-file=diff.txt --output-file=summary.txt --print`,
		Run:     summaryGitDiff,
		Args:    cobra.MaximumNArgs(1),
	}
)

// flags
var (
	// inputFile is the input file that contains the git diff.
	inputFile string
	// outputFile is the output file that contains the summary of the git diff.
	outputFile string
	// print is a flag to print the summary of the git diff to the console.
	print bool
	// aiModel is the model to use for the AI.
	aiModel string
	// aiProvider is the provider to use for the AI.
	aiProvider string
	// customPrompt is the custom prompt message for the AI.
	customPrompt string
)

func init() {
	SummaryGitDiff.PersistentFlags().StringVar(&inputFile, "input-file", "", "The input file that contains the git diff")
	SummaryGitDiff.PersistentFlags().StringVar(&outputFile, "output-file", "", "The output file that contains the summary of the git diff")
	SummaryGitDiff.PersistentFlags().BoolVar(&print, "print", false, "Print the summary of the git diff to the console")
	SummaryGitDiff.PersistentFlags().StringVar(&aiModel, "ai-model", "", "The model to use for the AI")
	SummaryGitDiff.PersistentFlags().StringVar(&aiProvider, "ai-provider", "openai", "The provider to use for the AI")
	SummaryGitDiff.PersistentFlags().StringVar(&customPrompt, "custom-prompt", "", "The custom prompt message for the AI")
}

var (
	gitDiffSummaryPrompt = map[ai.ProviderKind]string{
		ai.OpenAIProvider: `You are expert programmer and you receive multiple of git diff, please went over every diff file and created a concise summary for each one. Use the following format of conventional-commits and writing rule to summarize it:
format of conventional-commits:
"### dir/file_name
<type>[optional scope]: <Short description, up to 15 chars>

<detailed description text> - Be specific and comprehensive. Wrap it to 30 characters."

writing rules:
- Use the imperative, present tense: "change" not "changed" nor "changes".
- Don't capitalize first letter of detailed description and short description.
- No dot (.) at the end.
- Don't put actual diff content into the summary.
	
git diff content:
\n\n
`,
	}
)

func summaryGitDiff(cmd *cobra.Command, args []string) {
	diffContent := []byte{}
	if len(args) == 1 {
		diffContent = []byte(args[0])
	} else {
		diffContentBytes, err := os.ReadFile(inputFile)
		if err != nil {
			log.Fatalf("summary-gitdiff: read git diff content from %s, failed, err:%+v", inputFile, err)
		}
		diffContent = diffContentBytes
	}
	diffFiles, err := git.ParseDiffContentToMultiFileDiff(diffContent)
	if err != nil {
		log.Fatalf("summary-gitdiff: parse git diff content failed, err:%+v, please check git diff content", err)
	}

	payloads := git.MultiFileDiffToPayloadSlice(diffFiles, 30000)

	aiAdapter, err := NewAIAdapter(ai.ProviderKind(aiProvider), aiModel)
	if err != nil {
		log.Fatalf("summary-gitdiff: create AI adapter failed, err:%+v", err)
	}

	promptStr, ok := gitDiffSummaryPrompt[aiAdapter.GetProvider()]
	if !ok {
		log.Fatalf("summary-gitdiff: provider %s is not supported summary prompt message", aiAdapter.GetProvider())
	}
	if customPrompt != "" {
		promptStr = customPrompt
	}
	summaryResult := ""
	for _, payload := range payloads {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		summary, err := aiAdapter.SimplePrompt(ctx, aiModel, promptStr+payload)
		if err != nil {
			cancel()
			log.Fatalf("summary-gitdiff: prompt AI model failed, err:%+v", err)
		}
		cancel()
		if summaryResult != "" {
			summaryResult += "\n"
		}

		summaryResult += summary
	}
	if print {
		fmt.Println(summaryResult)
	}
	if outputFile != "" {
		err := os.WriteFile(outputFile, []byte(summaryResult), 0644)
		if err != nil {
			log.Fatalf("summary-gitdiff: write summary to %s failed, err:%+v", outputFile, err)
		}
	}
}
