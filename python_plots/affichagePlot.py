import matplotlib.pyplot as plt
import numpy as np

# Données du champ PersonnalityAverage
personality_average_data = ...
nbSimulations = ...

categories = list(personality_average_data.keys())
subcategories = list(personality_average_data[categories[0]].keys())

# Création de 4 plots distincts
fig, axes = plt.subplots(nrows=2, ncols=2, figsize=(10, 8))
fig.suptitle('Personality Traits Average Values')

for i, ax in enumerate(axes.flat):
    category = categories[i]
    values = [personality_average_data[category][subcat]/nbSimulations for subcat in subcategories]
    ax.plot(subcategories, values, alpha=0.7)
    ax.set_title(category)
    ax.set_xlabel('Subcategories')
    ax.set_ylabel('Average Values')

plt.tight_layout(rect=[0, 0.03, 1, 0.95])
plt.show()