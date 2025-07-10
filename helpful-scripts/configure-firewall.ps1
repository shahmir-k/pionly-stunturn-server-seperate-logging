# WebRTC Server Firewall Configuration Script
# This script configures Windows Firewall to allow the WebRTC server to function properly

Write-Host "Configuring Windows Firewall for WebRTC Server..." -ForegroundColor Green
Write-Host ""

# Check if running as administrator
if (-NOT ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
    Write-Host "ERROR: This script must be run as Administrator!" -ForegroundColor Red
    Write-Host "Please right-click PowerShell and select 'Run as Administrator'" -ForegroundColor Yellow
    pause
    exit 1
}

# Function to create firewall rule
function New-FirewallRule {
    param(
        [string]$Name,
        [string]$Direction,
        [string]$Protocol,
        [string]$LocalPort,
        [string]$Action = "Allow"
    )
    
    try {
        # Check if rule already exists
        $existingRule = Get-NetFirewallRule -DisplayName $Name -ErrorAction SilentlyContinue
        
        if ($existingRule) {
            Write-Host "Rule '$Name' already exists. Updating..." -ForegroundColor Yellow
            Remove-NetFirewallRule -DisplayName $Name -Confirm:$false
        }
        
        # Create new rule
        New-NetFirewallRule -DisplayName $Name `
                          -Direction $Direction `
                          -Protocol $Protocol `
                          -LocalPort $LocalPort `
                          -Action $Action `
                          -Profile Any `
                          -Description "WebRTC Server - $Name"
        
        Write-Host "✓ Created firewall rule: $Name" -ForegroundColor Green
    }
    catch {
        Write-Host "✗ Failed to create rule '$Name': $($_.Exception.Message)" -ForegroundColor Red
    }
}

# Create firewall rules for WebRTC server

Write-Host "Creating inbound rules..." -ForegroundColor Cyan

# HTTPS/WebSocket Signaling Server (Port 443)
New-FirewallRule -Name "WebRTC Signaling HTTPS" -Direction "Inbound" -Protocol "TCP" -LocalPort "443"

# STUN/TURN UDP (Port 3478)
New-FirewallRule -Name "WebRTC STUN/TURN UDP" -Direction "Inbound" -Protocol "UDP" -LocalPort "3478"

# STUN/TURN TCP (Port 3478)
New-FirewallRule -Name "WebRTC STUN/TURN TCP" -Direction "Inbound" -Protocol "TCP" -LocalPort "3478"

# STUN/TURN TLS (Port 5349) - TLS only uses TCP
New-FirewallRule -Name "WebRTC STUN/TURN TLS" -Direction "Inbound" -Protocol "TCP" -LocalPort "5349"

Write-Host ""
Write-Host "Creating outbound rules..." -ForegroundColor Cyan

# Outbound rules for the same ports (in case of symmetric NAT)
New-FirewallRule -Name "WebRTC Signaling HTTPS Outbound" -Direction "Outbound" -Protocol "TCP" -LocalPort "443"
New-FirewallRule -Name "WebRTC STUN/TURN UDP Outbound" -Direction "Outbound" -Protocol "UDP" -LocalPort "3478"
New-FirewallRule -Name "WebRTC STUN/TURN TCP Outbound" -Direction "Outbound" -Protocol "TCP" -LocalPort "3478"
New-FirewallRule -Name "WebRTC STUN/TURN TLS Outbound" -Direction "Outbound" -Protocol "TCP" -LocalPort "5349"

Write-Host ""
Write-Host "Firewall configuration completed!" -ForegroundColor Green
Write-Host ""
Write-Host "Summary of configured ports:" -ForegroundColor Yellow
Write-Host "• Port 443 (TCP) - HTTPS/WebSocket Signaling Server" -ForegroundColor White
Write-Host "• Port 3478 (UDP/TCP) - STUN/TURN Server" -ForegroundColor White
Write-Host "• Port 5349 (TCP) - STUN/TURN TLS Server" -ForegroundColor White
Write-Host ""
Write-Host "You can now run your WebRTC server without firewall issues." -ForegroundColor Green
Write-Host "To remove these rules later, use: Remove-NetFirewallRule -DisplayName WebRTC" -ForegroundColor Gray

pause 