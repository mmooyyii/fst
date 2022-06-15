def draw_graph(edges):
    import networkx as nx
    import matplotlib.pyplot as plt
    edges = [list(map(str, edge)) for edge in edges]
    edge_labels = {}
    G = nx.DiGraph()
    for i, v in enumerate(edges):
        a, b, label = v
        G.add_edge(a, b, id=i + 1)
        edge_labels[(a, b)] = label
    pos = nx.spring_layout(G)
    nx.draw_networkx_nodes(G, pos, cmap=plt.get_cmap('jet'), node_size=500)
    nx.draw(
        G, pos, edge_color='black', width=1, linewidths=1,
        node_size=500, node_color='pink', alpha=0.9,
        labels={node: node for node in G.nodes()}
    )
    nx.draw_networkx_edge_labels(
        G, pos, font_color='red',
        edge_labels=edge_labels
    )
    plt.show()


import json

graph = input()

graph = [[i['a'], i['b'], i['char'] + ':' + i['output']] for i in json.loads(graph)]
draw_graph(graph)
