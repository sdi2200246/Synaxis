import numpy as np
from collections import defaultdict
from math import log2

def rmse(model, ratings):
    preds = []
    truth = []
    for u, e, r in ratings:
        if u in model.user_index and e in model.event_index:
            preds.append(model.predict(u, e))
            truth.append(r)
    preds = np.array(preds, dtype=float)
    truth = np.array(truth, dtype=float)
    return float(np.sqrt(np.mean((truth - preds) ** 2))) if len(truth) else float("nan")


def mae(model, ratings):
    preds = []
    truth = []
    for u, e, r in ratings:
        if u in model.user_index and e in model.event_index:
            preds.append(model.predict(u, e))
            truth.append(r)
    preds = np.array(preds, dtype=float)
    truth = np.array(truth, dtype=float)
    return float(np.mean(np.abs(truth - preds))) if len(truth) else float("nan")


def ndcg_at_k(ranked_items, relevant_items, k=10):
    """
    ranked_items: list of event_ids in predicted order
    relevant_items: set of relevant event_ids
    """
    dcg = 0.0
    for idx, item in enumerate(ranked_items[:k], start=1):
        if item in relevant_items:
            dcg += 1.0 / log2(idx + 1)

    ideal_hits = min(len(relevant_items), k)
    idcg = sum(1.0 / log2(i + 1) for i in range(1, ideal_hits + 1))

    return dcg / idcg if idcg > 0 else 0.0


def precision_recall_hit_ndcg_at_k(model, train_ratings, test_ratings, k=10, relevance_threshold=4.0):
    """
    Evaluates ranking quality on held-out test ratings.

    relevance_threshold:
        ratings >= threshold are considered relevant.
        In your data, test bookings are 4.0, so 4.0 is a good choice.
    """
    train_seen = defaultdict(set)
    test_relevant = defaultdict(set)

    for u, e, r in train_ratings:
        train_seen[u].add(e)

    for u, e, r in test_ratings:
        if r >= relevance_threshold:
            test_relevant[u].add(e)

    precisions = []
    recalls = []
    hit_rates = []
    ndcgs = []

    for u, relevant_items in test_relevant.items():
        if u not in model.user_index:
            continue

        seen = train_seen.get(u, set())

        # rank all candidate events except ones already seen in training
        scores = []
        for event_id in model.event_index.keys():
            if event_id in seen:
                continue
            score = model.predict(u, event_id)
            scores.append((event_id, score))

        scores.sort(key=lambda x: x[1], reverse=True)
        top_k = [event_id for event_id, _ in scores[:k]]

        hits = sum(1 for item in top_k if item in relevant_items)

        precision = hits / k
        recall = hits / len(relevant_items) if relevant_items else 0.0
        hit_rate = 1.0 if hits > 0 else 0.0
        ndcg = ndcg_at_k(top_k, relevant_items, k=k)

        precisions.append(precision)
        recalls.append(recall)
        hit_rates.append(hit_rate)
        ndcgs.append(ndcg)

    return {
        f"Precision@{k}": float(np.mean(precisions)) if precisions else 0.0,
        f"Recall@{k}": float(np.mean(recalls)) if recalls else 0.0,
        f"HitRate@{k}": float(np.mean(hit_rates)) if hit_rates else 0.0,
        f"NDCG@{k}": float(np.mean(ndcgs)) if ndcgs else 0.0,
    }


def evaluate_model(model, train_ratings, test_ratings, k=10, relevance_threshold=4.0):
    rating_rmse = rmse(model, test_ratings)
    rating_mae = mae(model, test_ratings)
    ranking_metrics = precision_recall_hit_ndcg_at_k(
        model,
        train_ratings=train_ratings,
        test_ratings=test_ratings,
        k=k,
        relevance_threshold=relevance_threshold,
    )

    results = {
        "RMSE": rating_rmse,
        "MAE": rating_mae,
        **ranking_metrics,
    }
    return results