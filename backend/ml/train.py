from dataloader import FileDataLoader , DataLoader
from model import BiasedMF
from eval import evaluate_model
import matplotlib.pyplot as plt
 
 
def train(loader : DataLoader ,k =10, alpha=0.00001, lam=0.5, epochs=200):
    
    ratings = loader.load_ratings()
    
    print("── merged ratings (visit + booking) ─────────────────")
    for u, e, r in sorted(ratings)[:10]:
        print(f"  {u}  {e}  {r:.1f}")
 
    model = BiasedMF(k=k, alpha=alpha, lam=lam, epochs=epochs)
    history = model.fit(ratings)
    print(f"\nfinal loss: {history[-1]:.6f}")

    plt.figure(figsize=(8,5))
    plt.plot(history, linewidth=2)
    plt.title("Training Loss")
    plt.xlabel("Epoch")
    plt.ylabel("Loss")
    plt.grid(True, alpha=0.3)
    plt.show()

    return model, ratings
    
 
 
if __name__ == "__main__":
    train_loader = FileDataLoader("train_recs_medium.json")
    test_loader  = FileDataLoader("test_recs_medium.json")

    train_ratings = train_loader.load_ratings()
    test_ratings = test_loader.load_ratings()

    model, _ = train(loader=train_loader,k=15,alpha=0.01,lam=0.05,epochs=200)

    results = evaluate_model(model, train_ratings, test_ratings, k=5, relevance_threshold=4.0)
    print("\nEvaluation results:")
    for name, value in results.items():
        print(f"{name}: {value:.4f}")