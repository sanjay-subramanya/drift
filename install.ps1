$repo = "sanjay-subramanya/drift"
$binaryName = "drift"

# 1. Get Release Data
$release = Invoke-RestMethod -Uri "https://api.github.com/repos/$repo/releases/latest"
$asset = $release.assets | Where-Object { $_.name -like "*windows_amd64.zip" } | Select-Object -First 1

if ($null -eq $asset) { Write-Error "No Windows build found."; exit }

# 2. Download and Extract
Write-Host "Downloading $($asset.name)..."
Invoke-WebRequest -Uri $asset.browser_download_url -OutFile "drift_dist.zip"
Expand-Archive -Path "drift_dist.zip" -DestinationPath "." -Force
Remove-Item "drift_dist.zip"

Write-Host "Installation Complete. Move drift.exe to a folder in your PATH."