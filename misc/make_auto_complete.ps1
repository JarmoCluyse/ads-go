Register-ArgumentCompleter -Native -CommandName make -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)

    # Check if a Makefile exists in the current directory
    if (Test-Path .\Makefile) {
        # Parse targets (lines starting with a word followed by a colon)
        # Filters out comments, internal targets (starting with .), and variables
        Get-Content .\Makefile | 
            Where-Object { $_ -match '^[a-zA-Z0-9_-]+\s*:' } | 
            ForEach-Object { $_.Split(':')[0].Trim() } | 
            Where-Object { $_ -like "$wordToComplete*" } | 
            ForEach-Object { 
                [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_) 
            }
    }
}