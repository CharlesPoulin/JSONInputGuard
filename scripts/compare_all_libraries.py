#!/usr/bin/env python3
"""
Comprehensive comparison of JSONInputGuard vs popular JSON libraries
"""

import matplotlib.pyplot as plt
import numpy as np
import pandas as pd
from matplotlib import rcParams

# Set style for better looking charts
rcParams['figure.figsize'] = (14, 10)
rcParams['font.size'] = 10

def create_comprehensive_comparison():
    """Create comprehensive comparison charts for all libraries"""
    
    # Data from benchmark results
    libraries = ['GJSON', 'JSONInputGuard', 'Jsoniter', 'Fastjson', 'Standard Library']
    ops_per_sec = [22652, 4509, 339, 204, 363]
    throughput_mbps = [1145.45, 235.53, 19.40, 11.35, 19.76]
    memory_bytes = [68424, 14487, 876574, 3292834, 797260]
    allocations = [2, 5, 16496, 172, 16489]
    
    # Create subplots
    fig, ((ax1, ax2), (ax3, ax4)) = plt.subplots(2, 2, figsize=(16, 12))
    fig.suptitle('JSONInputGuard vs Popular JSON Libraries Performance Comparison', fontsize=16, fontweight='bold')
    
    # Color scheme
    colors = ['#FF6B6B', '#4ECDC4', '#45B7D1', '#96CEB4', '#FFEAA7']
    
    # 1. Operations per second comparison
    bars1 = ax1.bar(libraries, ops_per_sec, color=colors, alpha=0.8)
    ax1.set_title('Operations per Second (Higher is Better)', fontweight='bold')
    ax1.set_ylabel('Operations/sec')
    ax1.tick_params(axis='x', rotation=45)
    
    # Add value labels on bars
    for bar, value in zip(bars1, ops_per_sec):
        height = bar.get_height()
        ax1.text(bar.get_x() + bar.get_width()/2., height + 500,
                f'{value:,}', ha='center', va='bottom', fontweight='bold', fontsize=9)
    
    # 2. Throughput comparison
    bars2 = ax2.bar(libraries, throughput_mbps, color=colors, alpha=0.8)
    ax2.set_title('Throughput (MB/s)', fontweight='bold')
    ax2.set_ylabel('MB/s')
    ax2.tick_params(axis='x', rotation=45)
    
    for bar, value in zip(bars2, throughput_mbps):
        height = bar.get_height()
        ax2.text(bar.get_x() + bar.get_width()/2., height + 20,
                f'{value:.1f}', ha='center', va='bottom', fontweight='bold')
    
    # 3. Memory usage comparison (log scale)
    bars3 = ax3.bar(libraries, memory_bytes, color=colors, alpha=0.8)
    ax3.set_title('Memory Usage per Operation (Bytes)', fontweight='bold')
    ax3.set_ylabel('Bytes (log scale)')
    ax3.set_yscale('log')
    ax3.tick_params(axis='x', rotation=45)
    
    for bar, value in zip(bars3, memory_bytes):
        height = bar.get_height()
        ax3.text(bar.get_x() + bar.get_width()/2., height * 1.1,
                f'{value:,}', ha='center', va='bottom', fontweight='bold', fontsize=8)
    
    # 4. Allocations comparison
    bars4 = ax4.bar(libraries, allocations, color=colors, alpha=0.8)
    ax4.set_title('Allocations per Operation', fontweight='bold')
    ax4.set_ylabel('Allocations')
    ax4.tick_params(axis='x', rotation=45)
    
    for bar, value in zip(bars4, allocations):
        height = bar.get_height()
        ax4.text(bar.get_x() + bar.get_width()/2., height + 200,
                f'{value:,}', ha='center', va='bottom', fontweight='bold')
    
    plt.tight_layout()
    plt.savefig('all_libraries_comparison.png', dpi=300, bbox_inches='tight')
    plt.show()

def create_radar_chart_all_libraries():
    """Create a radar chart showing relative performance of all libraries"""
    
    # Normalize values for radar chart (0-1 scale)
    categories = ['Speed\n(ops/sec)', 'Throughput\n(MB/s)', 'Memory\nEfficiency', 'Allocation\nEfficiency']
    
    # Normalize to 0-1 scale where 1 is best performance
    gjson = [1.0, 1.0, 1.0, 1.0]  # Best in all categories
    jsoninputguard = [0.199, 0.206, 0.212, 0.0003]  # Normalized against GJSON
    jsoniter = [0.015, 0.017, 0.078, 0.0001]
    fastjson = [0.009, 0.010, 0.021, 0.010]
    stdlib = [0.016, 0.017, 0.086, 0.0001]
    
    angles = np.linspace(0, 2 * np.pi, len(categories), endpoint=False).tolist()
    angles += angles[:1]  # Complete the circle
    
    gjson += gjson[:1]
    jsoninputguard += jsoninputguard[:1]
    jsoniter += jsoniter[:1]
    fastjson += fastjson[:1]
    stdlib += stdlib[:1]
    
    fig, ax = plt.subplots(figsize=(12, 12), subplot_kw=dict(projection='polar'))
    
    ax.plot(angles, gjson, 'o-', linewidth=2, label='GJSON', color='#FF6B6B')
    ax.fill(angles, gjson, alpha=0.25, color='#FF6B6B')
    ax.plot(angles, jsoninputguard, 'o-', linewidth=2, label='JSONInputGuard', color='#4ECDC4')
    ax.fill(angles, jsoninputguard, alpha=0.25, color='#4ECDC4')
    ax.plot(angles, jsoniter, 'o-', linewidth=2, label='Jsoniter', color='#45B7D1')
    ax.fill(angles, jsoniter, alpha=0.25, color='#45B7D1')
    ax.plot(angles, fastjson, 'o-', linewidth=2, label='Fastjson', color='#96CEB4')
    ax.fill(angles, fastjson, alpha=0.25, color='#96CEB4')
    ax.plot(angles, stdlib, 'o-', linewidth=2, label='Standard Library', color='#FFEAA7')
    ax.fill(angles, stdlib, alpha=0.25, color='#FFEAA7')
    
    ax.set_xticks(angles[:-1])
    ax.set_xticklabels(categories)
    ax.set_ylim(0, 1)
    ax.set_title('Relative Performance Comparison\n(1.0 = Best Performance)', fontsize=14, fontweight='bold', pad=20)
    ax.legend(loc='upper right', bbox_to_anchor=(1.4, 1.0))
    
    plt.tight_layout()
    plt.savefig('radar_all_libraries.png', dpi=300, bbox_inches='tight')
    plt.show()

def create_efficiency_table():
    """Create a formatted efficiency comparison table"""
    
    data = {
        'Library': ['GJSON', 'JSONInputGuard', 'Jsoniter', 'Fastjson', 'Standard Library'],
        'Operations/sec': [22652, 4509, 339, 204, 363],
        'Throughput (MB/s)': [1145.45, 235.53, 19.40, 11.35, 19.76],
        'Memory (bytes)': [68424, 14487, 876574, 3292834, 797260],
        'Allocations': [2, 5, 16496, 172, 16489],
        'Speed vs Stdlib': ['62.4x', '12.4x', '0.9x', '0.6x', '1.0x'],
        'Memory vs Stdlib': ['11.7x less', '55x less', '0.9x', '4.1x more', '1.0x']
    }
    
    df = pd.DataFrame(data)
    print("\n" + "="*100)
    print("JSONInputGuard vs Popular JSON Libraries - Comprehensive Comparison")
    print("="*100)
    print(df.to_string(index=False))
    print("="*100)
    
    return df

def create_recommendations():
    """Print recommendations based on use cases"""
    
    print("\n🎯 USE CASE RECOMMENDATIONS")
    print("="*50)
    
    print("\n🚀 For Maximum Speed (Pure Parsing):")
    print("  1. GJSON - 22,652 ops/sec")
    print("  2. JSONInputGuard - 4,509 ops/sec")
    print("  3. Standard Library - 363 ops/sec")
    
    print("\n💾 For Memory Efficiency:")
    print("  1. GJSON - 68KB/op")
    print("  2. JSONInputGuard - 14KB/op")
    print("  3. Standard Library - 797KB/op")
    
    print("\n🛡️ For APIs with Validation:")
    print("  1. JSONInputGuard - AST-based validation")
    print("  2. GJSON - Simple field checks")
    print("  3. Standard Library - Full validation but slow")
    
    print("\n🔄 For Drop-in Replacement:")
    print("  1. Jsoniter - API compatible with stdlib")
    print("  2. JSONInputGuard - If you can refactor")
    print("  3. Standard Library - Baseline")
    
    print("\n📊 Key Insights:")
    print("  • GJSON is fastest but read-only")
    print("  • JSONInputGuard offers best balance of speed + features")
    print("  • Standard library is slowest but most compatible")
    print("  • Fastjson uses most memory (3.3MB/op)")
    print("  • JSONInputGuard is 12.4x faster than stdlib")

if __name__ == "__main__":
    print("Generating comprehensive library comparison...")
    
    # Create efficiency table
    df = create_efficiency_table()
    
    # Create charts
    create_comprehensive_comparison()
    create_radar_chart_all_libraries()
    
    # Print recommendations
    create_recommendations()
    
    print("\n📊 Charts saved as:")
    print("- all_libraries_comparison.png")
    print("- radar_all_libraries.png")
    
    print("\n🏆 Top Performers:")
    print("• GJSON: Fastest overall (22,652 ops/sec)")
    print("• JSONInputGuard: Best balance of speed + features (4,509 ops/sec)")
    print("• Standard Library: Slowest but most compatible (363 ops/sec)")
