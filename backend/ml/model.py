import numpy as np


class BiasedMF:
    """
    Biased Matrix Factorization trained with SGD.

    Prediction:
        r_hat[u,i] = mu + b_u[u] + b_i[i] + P[u] @ Q[i]

    Supports:
        fit(ratings)
        predict(user_id, event_id)
        top_n(user_id, n=20, exclude=None)
    """

    def __init__(
        self,
        k: int = 4,
        alpha: float = 0.01,
        lam: float = 0.05,
        epochs: int = 50,
        seed: int = 42,
        shuffle: bool = True,
    ):
        self.k = k
        self.alpha = alpha
        self.lam = lam
        self.epochs = epochs
        self.seed = seed
        self.shuffle = shuffle

        self.mu = 0.0
        self.b_u = None
        self.b_i = None
        self.P = None
        self.Q = None

        self.user_index: dict[str, int] = {}
        self.event_index: dict[str, int] = {}
        self.users = []
        self.events = []

    def fit(self, ratings: list[tuple[str, str, float]]) -> list[float]:
        """
        Train model with SGD.
        Returns loss history.
        """
        rng = np.random.default_rng(self.seed)

        # Build user/item indexes
        self.users = list(dict.fromkeys(u for u, _, _ in ratings))
        self.events = list(dict.fromkeys(i for _, i, _ in ratings))

        self.user_index = {u: idx for idx, u in enumerate(self.users)}
        self.event_index = {i: idx for idx, i in enumerate(self.events)}

        n_users = len(self.users)
        n_items = len(self.events)

        data = []
        for u_id, i_id, r in ratings:
            u = self.user_index[u_id]
            i = self.event_index[i_id]
            data.append((u, i, float(r)))

        self.mu = np.mean([r for _, _, r in data])

        # Initialize parameters
        self.b_u = np.zeros(n_users)
        self.b_i = np.zeros(n_items)
        self.P = rng.normal(0, 0.01, (n_users, self.k))
        self.Q = rng.normal(0, 0.01, (n_items, self.k))

        history = []

        for epoch in range(self.epochs):
            if self.shuffle:
                rng.shuffle(data)

            loss = self._step(data)
            history.append(loss)

            if (epoch + 1) % 10 == 0:
                print(f"epoch {epoch+1:>4}/{self.epochs}  loss={loss:.6f}")

        return history


    def _step(self, data: list[tuple[int, int, float]]) -> float:
        """
        One SGD epoch over all ratings.
        Updates parameters and returns epoch loss.
        """
        sq_err = 0.0

        for u, i, r in data:
            pred = self.mu + self.b_u[u] + self.b_i[i] + self.P[u] @ self.Q[i]
            e = r - pred

            sq_err += e * e

            p_old = self.P[u].copy()
            q_old = self.Q[i].copy()

            # biases
            self.b_u[u] += self.alpha * (e - self.lam * self.b_u[u])
            self.b_i[i] += self.alpha * (e - self.lam * self.b_i[i])

            # latent vectors
            self.P[u] += self.alpha * (e * q_old - self.lam * p_old)
            self.Q[i] += self.alpha * (e * p_old - self.lam * q_old)

        return self._loss(sq_err)


    def _loss(self, sq_err: float) -> float:
        """
        Full objective:
        squared error + L2 regularization
        """
        reg = self.lam * (
            np.sum(self.b_u ** 2)
            + np.sum(self.b_i ** 2)
            + np.sum(self.P ** 2)
            + np.sum(self.Q ** 2)
        )

        return float(sq_err + reg)
    
    def predict(self, user_id: str, event_id: str) -> float:
        if user_id not in self.user_index:
            return float(self.mu)

        if event_id not in self.event_index:
            u = self.user_index[user_id]
            return float(self.mu + self.b_u[u])

        u = self.user_index[user_id]
        i = self.event_index[event_id]

        return float(
            self.mu
            + self.b_u[u]
            + self.b_i[i]
            + self.P[u] @ self.Q[i]
        )

    def top_n(self,user_id: str,n: int = 20,exclude: set[str] | None = None,) -> list[tuple[str, float]]:

        if user_id not in self.user_index:
            scores = [
                (event_id, float(self.mu + self.b_i[i]))
                for event_id, i in self.event_index.items()
            ]
            scores.sort(key=lambda x: x[1], reverse=True)
            return scores[:n]

        u = self.user_index[user_id]

        raw_scores = (
            self.mu
            + self.b_u[u]
            + self.b_i
            + self.P[u] @ self.Q.T
        )

        results = []
        for event_id, i in self.event_index.items():
            if exclude and event_id in exclude:
                continue
            results.append((event_id, float(raw_scores[i])))

        results.sort(key=lambda x: x[1], reverse=True)
        return results[:n]
    

    def recomendations(self , k :int , bookings:list[tuple[str , str , float] , clicks: list[tuple[str, str, float]]]) -> list[tuple[str , str , float]]:
        all_recomendations = []
        for u in self.users:
            previous_bookings = set([e for user , e , _ in bookings if  user == u])
            all_recomendations.extend((u, e, score) for e, score in self.top_n(u, k , previous_bookings))
            
        return all_recomendations