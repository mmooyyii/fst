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
        labels={node: "" for node in G.nodes()}
    )
    nx.draw_networkx_edge_labels(
        G, pos, font_color='red',
        edge_labels=edge_labels
    )
    plt.show()


graph = [{"from": "1", "to": "3", "char": 97, "output": 2}, {"from": "3", "to": "4", "char": 98, "output": 0},
         {"from": "4", "to": "5", "char": 99, "output": 0}, {"from": "3", "to": "6", "char": 100, "output": 0},
         {"from": "6", "to": "4", "char": 99, "output": 0}, {"from": "4", "to": "5", "char": 99, "output": 0},
         {"from": "1", "to": "8", "char": 122, "output": 2}]

graph = [[i['from'], i['to'], chr(i['char']) + ':' + str(i['output'])] for i in graph]
draw_graph(graph)
