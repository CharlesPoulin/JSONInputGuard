#!/usr/bin/env python3
"""
Performance comparison visualization for JSONInputGuard vs Standard Library
"""

import matplotlib.pyplot as plt
import numpy as np
import pandas as pd
from matplotlib import rcParams

# Set style for better looking charts
rcParams['figure.figsize'] = (12, 8)
rcParams['font.size'] = 10

def create_performance_comparison():
    """Create comprehensive performance comparison charts"""
    
    # Data from your benchmark results
    benchmarks = ['Guard Only', 'Stdlib Unmarshal', 'Decode+Validate', 'HTTP Predict']
    ops_per_sec = [6494, 450, 286, 288]
    throughput_mbps = [351.01, 23.63, 15.01, 14.14]
    memory_bytes = [10089, 762791, 507283, 507803]
    allocations = [4, 16474, 8323, 8339]
    
    # Create subplots
    fig, ((ax1, ax2), (ax3, ax4)) = plt.subplots(2, 2, figsize=(15, 12))
    fig.suptitle('JSONInputGuard vs Standard Library Performance Comparison', fontsize=16, fontweight='bold')
    
    # 1. Operations per second comparison
    colors = ['#2E8B57', '#DC143C', '#4169E1', '#FF8C00']
    bars1 = ax1.bar(benchmarks, ops_per_sec, color=colors, alpha=0.8)
    ax1.set_title('Operations per Second (Higher is Better)', fontweight='bold')
    ax1.set_ylabel('Operations/sec')
    ax1.tick_params(axis='x', rotation=45)
    
    # Add value labels on bars
    for bar, value in zip(bars1, ops_per_sec):
        height = bar.get_height()
        ax1.text(bar.get_x() + bar.get_width()/2., height + 100,
                f'{value:,}', ha='center', va='bottom', fontweight='bold')
    
    # 2. Throughput comparison
    bars2 = ax2.bar(benchmarks, throughput_mbps, color=colors, alpha=0.8)
    ax2.set_title('Throughput (MB/s)', fontweight='bold')
    ax2.set_ylabel('MB/s')
    ax2.tick_params(axis='x', rotation=45)
    
    for bar, value in zip(bars2, throughput_mbps):
        height = bar.get_height()
        ax2.text(bar.get_x() + bar.get_width()/2., height + 5,
                f'{value:.1f}', ha='center', va='bottom', fontweight='bold')
    
    # 3. Memory usage comparison (log scale)
    bars3 = ax3.bar(benchmarks, memory_bytes, color=colors, alpha=0.8)
    ax3.set_title('Memory Usage per Operation (Bytes)', fontweight='bold')
    ax3.set_ylabel('Bytes (log scale)')
    ax3.set_yscale('log')
    ax3.tick_params(axis='x', rotation=45)
    
    for bar, value in zip(bars3, memory_bytes):
        height = bar.get_height()
        ax3.text(bar.get_x() + bar.get_width()/2., height * 1.1,
                f'{value:,}', ha='center', va='bottom', fontweight='bold', fontsize=8)
    
    # 4. Allocations comparison
    bars4 = ax4.bar(benchmarks, allocations, color=colors, alpha=0.8)
    ax4.set_title('Allocations per Operation', fontweight='bold')
    ax4.set_ylabel('Allocations')
    ax4.tick_params(axis='x', rotation=45)
    
    for bar, value in zip(bars4, allocations):
        height = bar.get_height()
        ax4.text(bar.get_x() + bar.get_width()/2., height + 200,
                f'{value:,}', ha='center', va='bottom', fontweight='bold')
    
    plt.tight_layout()
    plt.savefig('performance_comparison.png', dpi=300, bbox_inches='tight')
    plt.show()

def create_radar_chart():
    """Create a radar chart showing relative performance"""
    
    # Normalize values for radar chart (0-1 scale)
    categories = ['Speed\n(ops/sec)', 'Throughput\n(MB/s)', 'Memory\nEfficiency', 'Allocation\nEfficiency']
    
    # Normalize to 0-1 scale where 1 is best performance
    guard_only = [1.0, 1.0, 1.0, 1.0]  # Best in all categories
    stdlib = [0.069, 0.067, 0.013, 0.0002]  # Normalized against guard
    
    angles = np.linspace(0, 2 * np.pi, len(categories), endpoint=False).tolist()
    angles += angles[:1]  # Complete the circle
    
    guard_only += guard_only[:1]
    stdlib += stdlib[:1]
    
    fig, ax = plt.subplots(figsize=(10, 10), subplot_kw=dict(projection='polar'))
    
    ax.plot(angles, guard_only, 'o-', linewidth=2, label='JSONInputGuard', color='#2E8B57')
    ax.fill(angles, guard_only, alpha=0.25, color='#2E8B57')
    ax.plot(angles, stdlib, 'o-', linewidth=2, label='Standard Library', color='#DC143C')
    ax.fill(angles, stdlib, alpha=0.25, color='#DC143C')
    
    ax.set_xticks(angles[:-1])
    ax.set_xticklabels(categories)
    ax.set_ylim(0, 1)
    ax.set_title('Relative Performance Comparison\n(1.0 = Best Performance)', fontsize=14, fontweight='bold', pad=20)
    ax.legend(loc='upper right', bbox_to_anchor=(1.3, 1.0))
    
    plt.tight_layout()
    plt.savefig('radar_comparison.png', dpi=300, bbox_inches='tight')
    plt.show()

def create_efficiency_table():
    """Create a formatted efficiency comparison table"""
    
    data = {
        'Metric': ['Operations/sec', 'Throughput (MB/s)', 'Memory (bytes)', 'Allocations'],
        'Guard Only': [6494, 351.01, 10089, 4],
        'Stdlib': [450, 23.63, 762791, 16474],
        'Improvement': ['14.4x', '14.9x', '75x less', '4118x fewer'],
        'Percentage': ['93% faster', '93% faster', '98.7% less', '99.98% fewer']
    }
    
    df = pd.DataFrame(data)
    print("\n" + "="*80)
    print("JSONInputGuard vs Standard Library Efficiency Comparison")
    print("="*80)
    print(df.to_string(index=False))
    print("="*80)
    
    return df

if __name__ == "__main__":
    print("Generating performance comparison visualizations...")
    
    # Create efficiency table
    df = create_efficiency_table()
    
    # Create charts
    create_performance_comparison()
    create_radar_chart()
    
    print("\nCharts saved as:")
    print("- performance_comparison.png")
    print("- radar_comparison.png")
    print("\nKey Insights:")
    print("• Guard Only is 14.4x faster than standard library")
    print("• Uses 98.7% less memory per operation")
    print("• 99.98% fewer allocations (4 vs 16,474)")
    print("• 14.9x higher throughput (351 vs 23.6 MB/s)")

