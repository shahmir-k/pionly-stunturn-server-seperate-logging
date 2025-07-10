$Host.UI.RawUI.WindowTitle = "Signaling Log Monitor"
$logFile = "..\signaling.log"
$lastLineCount = 0

while (-not (Test-Path "..\shutdown-signal.txt")) {
    if (Test-Path $logFile) {
        $currentLineCount = (Get-Content $logFile).Count
        if ($currentLineCount -gt $lastLineCount) {
            $newLines = Get-Content $logFile | Select-Object -Skip $lastLineCount
            $newLines | ForEach-Object { Write-Host $_ }
            $lastLineCount = $currentLineCount
        }
    }
    Start-Sleep -Seconds 1
}
exit