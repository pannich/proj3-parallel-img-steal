import pandas as pd
import matplotlib
matplotlib.use('Agg')
import matplotlib.pyplot as plt

def plot_speedups(filename, title, output_filename):
    data = pd.read_csv(filename)
    plt.figure(figsize=(10, 6))
    for dataset in data['Dataset'].unique():
        subset = data[(data['Dataset'] == dataset)]
        plt.plot(subset['Thread'], subset['SpeedUp'], marker='x', linestyle='-', label=f'{dataset}')

    plt.title(f'Speedup vs. Number of Threads for {title}')
    plt.xlabel('Number of Threads')
    plt.ylabel('Speedup (1 thread runtime / n threads runtime)')
    plt.legend(title='Dataset')
    plt.grid(True)
    plt.xticks(data['Thread'].unique())
    plt.savefig(output_filename)
    plt.show()


def combine_plot(filepaths, labels, title, output_filename):
    plt.figure(figsize=(10, 6))

    # Iterate over each dataset
    for filepath, label in zip(filepaths, labels):
        data = pd.read_csv(filepath)
        for dataset in data['Dataset'].unique():
            subset = data[data['Dataset'] == dataset]
            plt.plot(subset['Thread'], subset['SpeedUp'], marker='x', linestyle='-', label=f'{label} - {dataset}')

    plt.title(f'Speedup Graph for {title}')
    plt.xlabel('Number of Threads')
    plt.ylabel('Speedup (1 thread runtime / n threads runtime)')
    plt.legend(title='Dataset')
    plt.grid(True)

    # Assume all data files have the same thread counts
    common_data = pd.read_csv(filepaths[0])
    plt.xticks(common_data['Thread'].unique())

    plt.savefig(output_filename)
    plt.close()

try:
    filepaths1 = [
        './parslices/small/speed_summary.txt',
        './parslicesBSP/small/speed_summary.txt',
        './parslicesBSPOptimized/small/speed_summary.txt'
    ]
    filepaths2 = [
        './parslices/whitespace/speed_summary.txt',  # Path to first dataset
        './parslicesBSP/whitespace/speed_summary.txt',  # Path to second dataset
        './parslicesBSPOptimized/whitespace/speed_summary.txt'
    ]
    labels = [
        'Parslices Method',
        'BSPStealing Method',
        'BSPStealingOptimized Method'
    ]
    combine_plot(filepaths1, labels, "Comparison of Methods x2 tasks", "small.png")
    combine_plot(filepaths2, labels, "Comparison of Methods x2 tasks", "whitespace.png")
except Exception as e:
    print(f"Error: Could not generate plots due to {e}")


# Generate plots for opts methods
# filepath = './parslicesBSPOptimized/whitespace/speed_summary'
# plot_speedups(filepath+".txt", filepath.split("/")[1], filepath + ".png")
