import { useState, useEffect } from "react";
import { Link, useSearchParams } from "react-router-dom";
import { ShoppingCart, Search, Package } from "lucide-react";

const CATEGORIES = [
  "Tümü",
  "Elektronik",
  "Bilgisayar",
  "Aksesuar",
  "Tablet",
  "Drone",
  "Kamera",
  "Akıllı Ev",
];

export default function ProductList({ onCartAdded }) {
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [activeCategory, setActiveCategory] = useState("Tümü");
  const [addedIds, setAddedIds] = useState({});
  const [searchParams] = useSearchParams();

  useEffect(() => {
    const cat = searchParams.get("category");
    if (cat) setActiveCategory(cat);
  }, [searchParams]);

  useEffect(() => {
    const fetchProducts = async () => {
      setLoading(true);
      try {
        let url = "http://localhost:3001/api/products?";
        if (activeCategory !== "Tümü") url += `category=${activeCategory}&`;
        if (search) url += `search=${search}`;
        const res = await fetch(url);
        const data = await res.json();
        setProducts(data);
      } catch (err) {
        console.error(err);
      }
      setLoading(false);
    };

    fetchProducts();
  }, [activeCategory, search]);

  const addToCart = async (e, productId) => {
    e.preventDefault();
    e.stopPropagation();

    const product = products.find((p) => p.ID === productId);
    if (!product || product.stock <= 0) {
      return;
    }

    try {
      await fetch("http://localhost:3001/api/cart", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ product_id: productId, quantity: 1 }),
      });
      setAddedIds((prev) => ({ ...prev, [productId]: true }));
      setTimeout(
        () => setAddedIds((prev) => ({ ...prev, [productId]: false })),
        1500,
      );
      onCartAdded(1);
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <div className="max-w-6xl mx-auto px-4 py-8">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-800">Ürünler</h1>
        <p className="text-sm text-gray-500 mt-1">
          Kategoriye göre filtreleyip sepete ekleyebilirsin.
        </p>
      </div>

      <div className="relative mb-4">
        <Search
          size={16}
          className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400"
        />
        <input
          type="text"
          placeholder="Ürün ara..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="w-full bg-white border border-gray-300 rounded-md pl-9 pr-3 py-2 text-sm text-gray-700 placeholder-gray-400 focus:outline-none focus:border-blue-500"
        />
      </div>

      <div className="flex gap-2 flex-wrap mb-6">
        {CATEGORIES.map((cat) => (
          <button
            key={cat}
            onClick={() => setActiveCategory(cat)}
            className={`px-3 py-1 rounded-md text-sm border ${
              activeCategory === cat
                ? "bg-blue-600 text-white border-blue-600"
                : "bg-white text-gray-700 border-gray-300 hover:bg-gray-50"
            }`}
          >
            {cat}
          </button>
        ))}
      </div>

      {loading ? (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {[...Array(8)].map((_, i) => (
            <div
              key={i}
              className="rounded-md bg-gray-200 h-64 animate-pulse"
            />
          ))}
        </div>
      ) : products.length === 0 ? (
        <div className="text-center py-12 text-gray-500">
          <Package size={40} className="mx-auto mb-3" />
          <p>Ürün bulunamadı</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {products.map((product) => (
            <Link
              key={product.ID}
              to={`/product/${product.ID}`}
              className="bg-white border border-gray-200 rounded-md overflow-hidden hover:shadow-sm"
            >
              <div className="h-44 bg-gray-100">
                <img
                  src={product.image}
                  alt={product.name}
                  className="w-full h-full object-cover"
                />
              </div>

              <div className="p-4">
                <div className="flex items-center justify-between mb-2 text-xs text-gray-500">
                  <span className="px-2 py-0.5 bg-gray-100 rounded">
                    {product.category}
                  </span>
                </div>

                <h3 className="text-gray-800 font-semibold text-sm mb-1">
                  {product.name}
                </h3>
                <p className="text-gray-500 text-xs mb-3">
                  {product.description}
                </p>

                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-blue-700 font-bold text-base">
                      ₺
                      {product.price.toLocaleString("tr-TR", {
                        minimumFractionDigits: 2,
                      })}
                    </p>
                    <p className="text-gray-500 text-xs">
                      {product.stock} adet stokta
                    </p>
                  </div>
                  <button
                    onClick={(e) => addToCart(e, product.ID)}
                    disabled={product.stock <= 0}
                    className={`flex items-center gap-1.5 px-3 py-1.5 rounded-md text-xs font-medium ${
                      product.stock <= 0
                        ? "bg-gray-200 text-gray-500 cursor-not-allowed"
                        : addedIds[product.ID]
                          ? "bg-green-100 text-green-700"
                          : "bg-blue-600 text-white hover:bg-blue-700"
                    }`}
                  >
                    <ShoppingCart size={13} />
                    {product.stock <= 0
                      ? "Stokta Yok"
                      : addedIds[product.ID]
                        ? "Eklendi!"
                        : "Sepete Ekle"}
                  </button>
                </div>
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
