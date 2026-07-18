"use client";

import { useEffect, useState } from "react";
import { api, PlaceSummary, Post, PriceInsight, SafetyAssessment } from "../lib/api";

export default function Home() {
  // Application state
  const [places, setPlaces] = useState<PlaceSummary[]>([]);
  const [selectedPlace, setSelectedPlace] = useState<string>("");
  const [posts, setPosts] = useState<Post[]>([]);
  const [loadingPlaces, setLoadingPlaces] = useState(true);
  const [loadingPosts, setLoadingPosts] = useState(false);

  // Price checker form state
  const [priceVertical, setPriceVertical] = useState<"taxi" | "exchange" | "food" | "tour" | "street_retail">("taxi");
  const [priceItem, setPriceItem] = useState("Airport taxi to Old Quarter");
  const [priceAmount, setPriceAmount] = useState("380000");
  const [priceCurrency, setPriceCurrency] = useState("VND");
  const [priceUnit, setPriceUnit] = useState("trip");
  const [priceRegion, setPriceRegion] = useState("hanoi-soc-son");
  const [priceSegment, setPriceSegment] = useState<"budget" | "standard" | "premium" | "luxury" | "regulated">("standard");
  const [priceVenue, setPriceVenue] = useState("transport_vendor");
  const [priceContext, setPriceContext] = useState("metered");
  const [priceInsight, setPriceInsight] = useState<PriceInsight | null>(null);
  const [checkingPrice, setCheckingPrice] = useState(false);

  // Safety assessment form state
  const [safetyFacts, setSafetyFacts] = useState("driver_refuses_to_stop");
  const [threatIndicator, setThreatIndicator] = useState("");
  const [injuryIndicator, setInjuryIndicator] = useState("");
  const [confinementIndicator, setConfinementIndicator] = useState("door_locked");
  const [coercionIndicator, setCoercionIndicator] = useState("");
  const [abilityToLeave, setAbilityToLeave] = useState<boolean>(false);
  const [safetyResult, setSafetyResult] = useState<SafetyAssessment | null>(null);
  const [checkingSafety, setCheckingSafety] = useState(false);

  // Forum post creation state
  const [newPostTitle, setNewPostTitle] = useState("");
  const [newPostBody, setNewPostBody] = useState("");
  const [newPostType, setNewPostType] = useState("review");
  const [postingDraft, setPostingDraft] = useState(false);
  const [draftPostId, setDraftPostId] = useState<string | null>(null);

  // Load places and initial posts on mount
  useEffect(() => {
    fetchPlaces();
    fetchPosts();
  }, []);

  const fetchPlaces = async () => {
    try {
      setLoadingPlaces(true);
      const res = await api.getPlaces();
      setPlaces(res.items);
    } catch (err) {
      console.error("Failed to load places", err);
    } finally {
      setLoadingPlaces(false);
    }
  };

  const fetchPosts = async (placeId?: string) => {
    try {
      setLoadingPosts(true);
      const res = await api.getPosts({ place_id: placeId });
      setPosts(res.items);
    } catch (err) {
      console.error("Failed to load posts", err);
    } finally {
      setLoadingPosts(false);
    }
  };

  const handlePlaceSelect = (placeId: string) => {
    setSelectedPlace(placeId);
    fetchPosts(placeId || undefined);
  };

  const handlePriceCheck = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      setCheckingPrice(true);
      setPriceInsight(null);
      const res = await api.checkPrice({
        vertical: priceVertical,
        raw_item: priceItem,
        money: {
          amount_minor: priceAmount,
          currency: priceCurrency,
          exponent: 0,
        },
        unit: priceUnit,
        region_id: priceRegion,
        service_segment: priceSegment,
        venue_type: priceVenue,
        transaction_context: priceContext,
        observed_at: new Date().toISOString(),
        extraction_confidence: 1.0,
        user_confirmed: true,
      });
      setPriceInsight(res);
    } catch (err: any) {
      alert("Price check failed: " + err.message);
    } finally {
      setCheckingPrice(false);
    }
  };

  const handleSafetyAssess = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      setCheckingSafety(true);
      setSafetyResult(null);

      const observed = safetyFacts.split(",").map(f => f.trim()).filter(Boolean);
      const threats = threatIndicator ? [threatIndicator] : [];
      const injuries = injuryIndicator ? [injuryIndicator] : [];
      const confinements = confinementIndicator ? [confinementIndicator] : [];
      const coercions = coercionIndicator ? [coercionIndicator] : [];

      const res = await api.assessSafety({
        observed_facts: observed,
        threat_indicators: threats,
        injury_indicators: injuries,
        confinement_indicators: confinements,
        coercion_indicators: coercions,
        ability_to_leave: abilityToLeave,
      });
      setSafetyResult(res);
    } catch (err: any) {
      alert("Safety assessment failed: " + err.message);
    } finally {
      setCheckingSafety(false);
    }
  };

  const handleCreatePost = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newPostTitle || !newPostBody) {
      alert("Please fill in all post fields");
      return;
    }
    try {
      setPostingDraft(true);
      // Create draft
      const draft = await api.createDraft({
        post_type: newPostType,
        original_locale: "vi-VN",
        title: newPostTitle,
        body: newPostBody,
        place_ids: selectedPlace ? [selectedPlace] : undefined,
      });
      setDraftPostId(draft.post_id);
      
      // Auto-publish for horizontal integration ease
      await api.publishPost(draft.post_id);
      
      setNewPostTitle("");
      setNewPostBody("");
      setDraftPostId(null);
      alert("Community Warning Post published successfully!");
      fetchPosts(selectedPlace || undefined);
    } catch (err: any) {
      alert("Failed to publish post: " + err.message);
    } finally {
      setPostingDraft(false);
    }
  };

  return (
    <div className="layout-container">
      {/* Premium Header */}
      <header style={{ marginBottom: "40px", borderBottom: "1px solid var(--border-glass)", paddingBottom: "24px" }}>
        <h1 style={{ fontSize: "3rem" }}>TOURTECT</h1>
        <p style={{ fontSize: "1.2rem", color: "var(--accent-secondary)", fontWeight: 500 }}>
          Travel Price Transparency & Safety Engine
        </p>
      </header>

      {/* Main Grid */}
      <div className="grid-cols-2">
        {/* Left Side: Directory and Forums */}
        <div style={{ display: "flex", flexDirection: "column", gap: "24px" }}>
          
          {/* Place Selector panel */}
          <section className="glass-panel">
            <h2 style={{ marginBottom: "16px", display: "flex", alignItems: "center", gap: "10px" }}>
              📍 Local Place Directory
            </h2>
            <div className="form-group">
              <label>Select region context filter:</label>
              <select 
                className="form-select"
                value={selectedPlace}
                onChange={(e) => handlePlaceSelect(e.target.value)}
              >
                <option value="">All Regions / Places</option>
                {places.map((place) => (
                  <option key={place.place_id} value={place.place_id}>
                    {place.name} ({place.category}) — {place.address || place.region_id}
                  </option>
                ))}
              </select>
            </div>
            {loadingPlaces && <p>Loading places from Postgres...</p>}
          </section>

          {/* Incident Warnings Board */}
          <section className="glass-panel" style={{ flexGrow: 1 }}>
            <h2 style={{ marginBottom: "16px", color: "var(--accent-warning)" }}>
              🛡️ Community Warnings Feed
            </h2>
            
            {loadingPosts ? (
              <p>Fetching incident reports...</p>
            ) : posts.length === 0 ? (
              <p style={{ color: "var(--text-muted)", fontStyle: "italic" }}>No active warning warnings in this area.</p>
            ) : (
              <div style={{ display: "flex", flexDirection: "column", gap: "16px", maxHeight: "400px", overflowY: "auto", paddingRight: "8px" }}>
                {posts.map((post) => (
                  <div 
                    key={post.post_id} 
                    style={{ 
                      padding: "16px", 
                      borderRadius: "8px", 
                      backgroundColor: "rgba(255, 255, 255, 0.03)", 
                      borderLeft: `4px solid ${
                        post.post_type === "scam_report" ? "var(--accent-error)" : "var(--accent-primary)"
                      }` 
                    }}
                  >
                    <div style={{ display: "flex", justifyContent: "space-between", marginBottom: "6px" }}>
                      <span style={{ fontWeight: 600, fontSize: "1.05rem" }}>{post.title}</span>
                      <span style={{ 
                        fontSize: "0.8rem", 
                        padding: "2px 8px", 
                        borderRadius: "12px", 
                        backgroundColor: post.post_type === "scam_report" ? "rgba(239, 68, 68, 0.15)" : "rgba(139, 92, 246, 0.15)",
                        color: post.post_type === "scam_report" ? "var(--accent-error)" : "var(--accent-primary)"
                      }}>
                        {post.post_type}
                      </span>
                    </div>
                    <p style={{ fontSize: "0.95rem", color: "var(--text-secondary)", marginBottom: "8px" }}>{post.body}</p>
                    <div style={{ fontSize: "0.8rem", color: "var(--text-muted)" }}>
                      Reported {new Date(post.created_at).toLocaleDateString()} UTC • Evidence: {post.evidence_level}
                    </div>
                  </div>
                ))}
              </div>
            )}

            {/* Create warning form */}
            <form onSubmit={handleCreatePost} style={{ marginTop: "24px", paddingTop: "20px", borderTop: "1px solid var(--border-glass)" }}>
              <h3 style={{ fontSize: "1.1rem", marginBottom: "12px" }}>File Community Incident Warning</h3>
              <div className="form-group">
                <input 
                  type="text" 
                  className="form-input" 
                  placeholder="Incident Title (e.g. Taxi meter scam at terminal)"
                  value={newPostTitle}
                  onChange={(e) => setNewPostTitle(e.target.value)}
                />
              </div>
              <div className="form-group">
                <textarea 
                  className="form-textarea" 
                  rows={3} 
                  placeholder="Provide objective facts (e.g. vehicle license, overcharged amount, driver dispute behavior)..."
                  value={newPostBody}
                  onChange={(e) => setNewPostBody(e.target.value)}
                />
              </div>
              <div style={{ display: "flex", gap: "12px" }}>
                <select 
                  className="form-select" 
                  style={{ width: "150px" }}
                  value={newPostType}
                  onChange={(e) => setNewPostType(e.target.value)}
                >
                  <option value="review">Review</option>
                  <option value="scam_report">Scam Report</option>
                  <option value="price_report">Price Report</option>
                  <option value="tip">Tip</option>
                </select>
                <button type="submit" className="neon-button" style={{ flexGrow: 1 }} disabled={postingDraft}>
                  {postingDraft ? "Publishing..." : "Broadcast Alert"}
                </button>
              </div>
            </form>
          </section>
        </div>

        {/* Right Side: Price Checker and Safety Engine */}
        <div style={{ display: "flex", flexDirection: "column", gap: "24px" }}>
          
          {/* Price Checker */}
          <section className="glass-panel">
            <h2 style={{ marginBottom: "16px", display: "flex", alignItems: "center", gap: "10px", color: "var(--accent-secondary)" }}>
              📊 Realtime Price Scan
            </h2>
            <form onSubmit={handlePriceCheck}>
              <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "12px", marginBottom: "12px" }}>
                <div className="form-group" style={{ marginBottom: 0 }}>
                  <label>Vertical</label>
                  <select className="form-select" value={priceVertical} onChange={(e: any) => setPriceVertical(e.target.value)}>
                    <option value="taxi">Taxi / Transport</option>
                    <option value="food">Food / Restaurant</option>
                    <option value="tour">Tour / Attraction</option>
                    <option value="street_retail">Street Retail</option>
                  </select>
                </div>
                <div className="form-group" style={{ marginBottom: 0 }}>
                  <label>Observed Cost</label>
                  <input type="text" className="form-input" value={priceAmount} onChange={(e) => setPriceAmount(e.target.value)} />
                </div>
              </div>

              <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "12px", marginBottom: "12px" }}>
                <div className="form-group" style={{ marginBottom: 0 }}>
                  <label>Item Name</label>
                  <input type="text" className="form-input" value={priceItem} onChange={(e) => setPriceItem(e.target.value)} />
                </div>
                <div className="form-group" style={{ marginBottom: 0 }}>
                  <label>Unit</label>
                  <input type="text" className="form-input" value={priceUnit} onChange={(e) => setPriceUnit(e.target.value)} />
                </div>
              </div>

              <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "12px", marginBottom: "16px" }}>
                <div className="form-group" style={{ marginBottom: 0 }}>
                  <label>Region ID</label>
                  <input type="text" className="form-input" value={priceRegion} onChange={(e) => setPriceRegion(e.target.value)} />
                </div>
                <div className="form-group" style={{ marginBottom: 0 }}>
                  <label>Context</label>
                  <select className="form-select" value={priceContext} onChange={(e) => setPriceContext(e.target.value)}>
                    <option value="metered">Metered</option>
                    <option value="verbal_quote">Verbal Quote</option>
                    <option value="negotiated">Negotiated</option>
                    <option value="posted_price">Posted Price</option>
                  </select>
                </div>
              </div>

              <button type="submit" className="neon-button" style={{ width: "100%" }} disabled={checkingPrice}>
                {checkingPrice ? "Analyzing Cohort Snapshots..." : "Scan Price Fair-Value"}
              </button>
            </form>

            {/* Price Result display */}
            {priceInsight && (
              <div 
                style={{ 
                  marginTop: "20px", 
                  padding: "16px", 
                  borderRadius: "8px", 
                  border: "1px solid var(--border-glass)",
                  backgroundColor: priceInsight.alert_level === "high_risk" ? "rgba(239, 68, 68, 0.08)" : "rgba(255, 255, 255, 0.02)"
                }}
              >
                <div style={{ display: "flex", justifyContent: "space-between", marginBottom: "8px" }}>
                  <span style={{ fontWeight: 600 }}>SCAN RESULT:</span>
                  <span style={{ 
                    fontWeight: 700, 
                    color: priceInsight.alert_level === "typical" ? "var(--accent-success)" : 
                           priceInsight.alert_level === "elevated" ? "var(--accent-warning)" : "var(--accent-error)"
                  }}>
                    {priceInsight.alert_level.toUpperCase()}
                  </span>
                </div>
                <p style={{ fontSize: "0.9rem", color: "var(--text-secondary)", marginBottom: "6px" }}>
                  Deviation: {(priceInsight.deviation_ratio || 0) > 0 ? "+" : ""}{Math.round((priceInsight.deviation_ratio || 0) * 100)}% from cohort median.
                </p>
                {priceInsight.reference && (
                  <div style={{ fontSize: "0.85rem", color: "var(--text-muted)", display: "grid", gridTemplateColumns: "1fr 1fr", gap: "6px", margin: "10px 0" }}>
                    <div>P10 Min: {priceInsight.reference.p10_minor} {priceInsight.reference.currency}</div>
                    <div>Median (P50): {priceInsight.reference.p50_minor} {priceInsight.reference.currency}</div>
                    <div>P90 Limit: {priceInsight.reference.p90_minor} {priceInsight.reference.currency}</div>
                    <div>Cohort size: {priceInsight.reference.effective_sample_size}</div>
                  </div>
                )}
                <div style={{ fontSize: "0.8rem", color: "var(--text-muted)", borderTop: "1px solid var(--border-glass)", paddingTop: "8px", marginTop: "8px" }}>
                  Reasons: {priceInsight.reasons.join(", ") || "none"} • Dataset: {priceInsight.dataset_version}
                </div>
              </div>
            )}
          </section>

          {/* Safety Shield Engine */}
          <section className="glass-panel">
            <h2 style={{ marginBottom: "16px", display: "flex", alignItems: "center", gap: "10px", color: "var(--accent-error)" }}>
              🛡️ Rule-First Safety Shield
            </h2>
            <form onSubmit={handleSafetyAssess}>
              <div className="form-group">
                <label>Observed facts:</label>
                <input 
                  type="text" 
                  className="form-input" 
                  value={safetyFacts} 
                  onChange={(e) => setSafetyFacts(e.target.value)} 
                  placeholder="e.g. price_dispute, driver_refuses_to_stop" 
                />
              </div>

              <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "12px", marginBottom: "16px" }}>
                <div className="form-group" style={{ marginBottom: 0 }}>
                  <label>Confinement Indicator</label>
                  <select className="form-select" value={confinementIndicator} onChange={(e) => setConfinementIndicator(e.target.value)}>
                    <option value="">None</option>
                    <option value="door_locked">Doors Locked by Driver</option>
                    <option value="physical_blocking">Path Blocked</option>
                  </select>
                </div>
                <div className="form-group" style={{ marginBottom: 0 }}>
                  <label>Threat Indicator</label>
                  <select className="form-select" value={threatIndicator} onChange={(e) => setThreatIndicator(e.target.value)}>
                    <option value="">None</option>
                    <option value="weapon">Weapon Brandished</option>
                    <option value="physical_violence">Physical Scuffle</option>
                    <option value="verbal_shouting">Intimidation / Shouting</option>
                  </select>
                </div>
              </div>

              <div className="form-group" style={{ display: "flex", alignItems: "center", gap: "8px" }}>
                <input 
                  type="checkbox" 
                  id="leaveCheck"
                  checked={!abilityToLeave} 
                  onChange={(e) => setAbilityToLeave(!e.target.checked)} 
                  style={{ width: "18px", height: "18px" }}
                />
                <label htmlFor="leaveCheck" style={{ margin: 0, cursor: "pointer" }}>Cannot safely exit/leave vehicle</label>
              </div>

              <button type="submit" className="neon-button" style={{ width: "100%", background: "linear-gradient(135deg, var(--accent-error) 0%, hsl(340, 80%, 45%) 100%)", boxShadow: "0 0 15px hsla(350, 85%, 55%, 0.25)" }} disabled={checkingSafety}>
                {checkingSafety ? "Evaluating Threat Levels..." : "Evaluate Safety Shield Playbook"}
              </button>
            </form>

            {/* Safety assessment result */}
            {safetyResult && (
              <div 
                style={{ 
                  marginTop: "20px", 
                  padding: "16px", 
                  borderRadius: "8px", 
                  border: "1px solid var(--border-glass)",
                  backgroundColor: safetyResult.urgency === "critical" ? "rgba(239, 68, 68, 0.1)" : "rgba(255, 255, 255, 0.02)"
                }}
              >
                <div style={{ display: "flex", justifyContent: "space-between", marginBottom: "12px" }}>
                  <span style={{ fontWeight: 600 }}>DISPUTE URGENCY:</span>
                  <span style={{ 
                    fontWeight: 800, 
                    color: safetyResult.urgency === "critical" ? "var(--accent-error)" : 
                           safetyResult.urgency === "urgent" ? "var(--accent-warning)" : "var(--accent-success)"
                  }}>
                    {safetyResult.urgency.toUpperCase()}
                  </span>
                </div>
                
                <div style={{ marginBottom: "12px" }}>
                  <span style={{ display: "block", fontSize: "0.85rem", fontWeight: 600, color: "var(--text-muted)", marginBottom: "4px" }}>
                    RECOMMENDED PLAYBOOK ACTIONS:
                  </span>
                  <ul style={{ paddingLeft: "18px", fontSize: "0.9rem", color: "var(--text-secondary)" }}>
                    {safetyResult.safe_actions.map((act, idx) => (
                      <li key={idx} style={{ marginBottom: "4px" }}>{act}</li>
                    ))}
                  </ul>
                </div>

                {safetyResult.silent_mode_recommended && (
                  <div style={{ padding: "8px 12px", borderRadius: "4px", backgroundColor: "rgba(239, 68, 68, 0.15)", color: "var(--accent-error)", fontSize: "0.85rem", fontWeight: 600, marginBottom: "8px", textAlign: "center" }}>
                    ⚠️ SILENT MODE RECOMMENDED — Do not make loud calls. Keep device flat.
                  </div>
                )}
                
                <div style={{ fontSize: "0.8rem", color: "var(--text-muted)", borderTop: "1px solid var(--border-glass)", paddingTop: "8px" }}>
                  Codes: {safetyResult.approved_action_codes?.join(", ") || "none"} • DB Ver: {safetyResult.safety_directory_version}
                </div>
              </div>
            )}
          </section>

        </div>
      </div>
    </div>
  );
}
