import numpy as np
import matplotlib.pyplot as plt
from scipy import integrate

class MagicValueCalculator:
    """
    Dynamic decay integral model for calculating user magic value
    
    M(t) = ∫_0^t (αU(τ)e^(-λ_u τ) - βD(τ)e^(-λ_d τ))dτ + M_0
    
    Where:
    - U(τ) and D(τ) represent upload/download rates at time τ (MB/s)
    - α=0.8, β=0.5 are contribution penalty coefficients
    - λ_u=1.2×10^(-5) s^(-1), λ_d=2.4×10^(-5) s^(-1) are time decay factors
    - M_0=100 is the initial value
    """
    
    def __init__(self, alpha=0.8, beta=0.5, lambda_u=1.2e-5, lambda_d=2.4e-5, M0=100):
        """Initialize model parameters"""
        self.alpha = alpha            # Upload contribution coefficient
        self.beta = beta              # Download penalty coefficient
        self.lambda_u = lambda_u      # Upload time decay factor (s^-1)
        self.lambda_d = lambda_d      # Download time decay factor (s^-1)
        self.M0 = M0                  # Initial magic value
    
    def integrand(self, t, upload_rate, download_rate):
        """Calculate the integrand value at time t"""
        return self.alpha * upload_rate * np.exp(-self.lambda_u * t) - \
               self.beta * download_rate * np.exp(-self.lambda_d * t)
    
    def calculate_magic_value(self, time_points, upload_rates, download_rates):
        """
        Calculate magic value at given time points
        
        Parameters:
        - time_points: Time series (seconds)
        - upload_rates: Upload rates at corresponding time points (MB/s)
        - download_rates: Download rates at corresponding time points (MB/s)
        
        Returns:
        - Magic value series at corresponding time points
        """
        if len(time_points) != len(upload_rates) or len(time_points) != len(download_rates):
            raise ValueError("Time points, upload rates, and download rates must have the same length")
        
        # 计算每个时间点的被积函数值
        integrand_values = [self.integrand(t, u, d) 
                            for t, u, d in zip(time_points, upload_rates, download_rates)]
        
        # 使用梯形法则进行数值积分
        magic_values = integrate.cumulative_trapezoid(integrand_values, time_points, initial=0) + self.M0
        
        return magic_values
    
    def plot_results(self, time_points, magic_values, upload_rates=None, download_rates=None):
        """Plot magic value and optional upload/download rates over time"""
        fig, ax1 = plt.subplots(figsize=(10, 6))
        
        # Plot magic value
        ax1.plot(time_points/3600, magic_values, 'b-', linewidth=2, label='Magic Value M(t)')
        ax1.set_xlabel('Time (hours)')
        ax1.set_ylabel('Magic Value', color='b')
        ax1.tick_params(axis='y', labelcolor='b')
        ax1.grid(True, linestyle='--', alpha=0.7)
        
        # If upload/download rates are provided, plot on secondary y-axis
        if upload_rates is not None and download_rates is not None:
            ax2 = ax1.twinx()
            ax2.plot(time_points/3600, upload_rates, 'g-', alpha=0.7, label='Upload Rate U(t)')
            ax2.plot(time_points/3600, download_rates, 'r-', alpha=0.7, label='Download Rate D(t)')
            ax2.set_ylabel('Rate (MB/s)', color='k')
            ax2.tick_params(axis='y', labelcolor='k')
            
            # Combine legends from both axes
            lines1, labels1 = ax1.get_legend_handles_labels()
            lines2, labels2 = ax2.get_legend_handles_labels()
            ax2.legend(lines1 + lines2, labels1 + labels2, loc='upper left')
        else:
            ax1.legend(loc='best')
        
        plt.title('User Magic Value Over Time')
        plt.tight_layout()
        plt.show()


# Create a combined plot function
def plot_combined_results(calculator, time_points, scenarios):
    """
    Plot multiple scenarios in a single figure with subplots
    
    Parameters:
    - calculator: MagicValueCalculator instance
    - time_points: Time points array
    - scenarios: List of dictionaries containing scenario data and metadata
    """
    # 设置4K分辨率 (3840x2160)
    fig = plt.figure(figsize=(16, 12), dpi=240)
    
    for i, scenario in enumerate(scenarios, 1):
        ax1 = fig.add_subplot(2, 2, i)
        
        # Plot magic value
        ax1.plot(time_points/3600, scenario['magic_values'], 'b-', linewidth=2, label='Magic Value M(t)')
        ax1.set_xlabel('Time (hours)')
        ax1.set_ylabel('Magic Value', color='b')
        ax1.tick_params(axis='y', labelcolor='b')
        ax1.grid(True, linestyle='--', alpha=0.7)
        
        # Add second y-axis for rates
        ax2 = ax1.twinx()
        ax2.plot(time_points/3600, scenario['upload_rates'], 'g-', alpha=0.7, label='Upload Rate U(t)')
        ax2.plot(time_points/3600, scenario['download_rates'], 'r-', alpha=0.7, label='Download Rate D(t)')
        ax2.set_ylabel('Rate (MB/s)', color='k')
        
        # Combine legends
        lines1, labels1 = ax1.get_legend_handles_labels()
        lines2, labels2 = ax2.get_legend_handles_labels()
        ax1.legend(lines1 + lines2, labels1 + labels2, loc='upper left', fontsize=8)
        
        ax1.set_title(f"Scenario {i}: {scenario['title']}")
    
    plt.suptitle('Magic Value Calculation Model Comparison', fontsize=16)
    plt.tight_layout()
    plt.subplots_adjust(top=0.92)
    
    # 保存为4K分辨率的PNG文件
    plt.savefig('magic_value_analysis_4k.png', dpi=240, bbox_inches='tight')
    plt.show()

# Example usage
if __name__ == "__main__":
    # Create calculator instance
    calculator = MagicValueCalculator()
    
    # Create a 24-hour time series (seconds)
    hours = 24
    time_points = np.linspace(0, hours * 3600, 1000)
    
    # Scenario 1: Constant upload and download rates
    constant_upload = np.ones_like(time_points) * 0.5  # Constant 0.5 MB/s upload
    constant_download = np.ones_like(time_points) * 0.8  # Constant 0.8 MB/s download
    magic_values_constant = calculator.calculate_magic_value(
        time_points, constant_upload, constant_download)
    
    # Scenario 2: Varying upload and download patterns
    # Simulate different upload/download behaviors at different times of day
    upload_rates = 0.3 + 0.5 * np.sin(2 * np.pi * time_points / (24 * 3600))
    download_rates = 0.5 + 0.8 * np.sin(2 * np.pi * time_points / (12 * 3600) + np.pi/4)
    magic_values_varying = calculator.calculate_magic_value(
        time_points, upload_rates, download_rates)
    
    # Scenario 3: Burst upload and download behavior
    # Low speed upload/download most of the time, but with several burst activities
    base_upload = np.ones_like(time_points) * 0.2
    base_download = np.ones_like(time_points) * 0.3
    
    # Add several upload bursts
    for peak_time in [3, 8, 15, 20]:  # Burst uploads at these hours
        peak_idx = np.abs(time_points - peak_time * 3600).argmin()
        window = 100  # Duration of burst in time points
        base_upload[peak_idx:peak_idx+window] += 2.0 * np.exp(-0.5 * 
            (np.arange(window) / (window/5))**2)
    
    # Add several download bursts
    for peak_time in [5, 10, 18, 22]:  # Burst downloads at these hours
        peak_idx = np.abs(time_points - peak_time * 3600).argmin()
        window = 80  # Duration of burst in time points
        base_download[peak_idx:peak_idx+window] += 1.5 * np.exp(-0.5 * 
            (np.arange(window) / (window/4))**2)
    
    magic_values_burst = calculator.calculate_magic_value(
        time_points, base_upload, base_download)
    
    # Scenario 4: Asymmetric usage pattern (new scenario)
    # Heavy uploads in first half of day, heavy downloads in second half
    asymmetric_upload = 0.8 * np.exp(-0.5 * ((time_points - 6*3600)/(3*3600))**2) + 0.2
    asymmetric_download = 0.8 * np.exp(-0.5 * ((time_points - 18*3600)/(4*3600))**2) + 0.3
    magic_values_asymmetric = calculator.calculate_magic_value(
        time_points, asymmetric_upload, asymmetric_download)
    
    # Collect all scenarios into a list for plotting
    scenarios = [
        {
            'title': 'Constant Rates',
            'upload_rates': constant_upload,
            'download_rates': constant_download,
            'magic_values': magic_values_constant
        },
        {
            'title': 'Periodic Variation',
            'upload_rates': upload_rates,
            'download_rates': download_rates,
            'magic_values': magic_values_varying
        },
        {
            'title': 'Burst Activity',
            'upload_rates': base_upload,
            'download_rates': base_download,
            'magic_values': magic_values_burst
        },
        {
            'title': 'Asymmetric Usage',
            'upload_rates': asymmetric_upload,
            'download_rates': asymmetric_download,
            'magic_values': magic_values_asymmetric
        }
    ]
    
    # Plot all scenarios in a single figure with subplots
    plot_combined_results(calculator, time_points, scenarios)