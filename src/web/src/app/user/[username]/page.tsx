"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { ArrowLeft, Lightbulb, FileText, Trophy, User as UserIcon } from "lucide-react";
import { getUserProfile, type UserProfile } from "@/lib/api";
import ErrorState from "@/components/ErrorState";

export default function UserProfilePage() {
  const { username } = useParams<{ username: string }>();
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!username) return;
    getUserProfile(username)
      .then(setProfile)
      .catch((err) => setError(err.message));
  }, [username]);

  const reload = () => {
    if (!username) return;
    setError(null);
    getUserProfile(username).then(setProfile).catch((err) => setError(err.message));
  };

  if (error) {
    return (
      <div className="mx-auto max-w-[860px] px-7 py-12">
        <ErrorState message={error} onRetry={reload} />
      </div>
    );
  }

  if (!profile) {
    return null; // loading.tsx handles this
  }

  const { user } = profile;

  return (
    <div className="mx-auto max-w-[860px] px-7 py-8">
      <Link href="/" className="mb-6 inline-flex items-center gap-1.5 text-sm text-ink-soft hover:text-ink">
        <ArrowLeft className="h-4 w-4" />
        返回
      </Link>

      {/* Profile card */}
      <div
        className="mb-8 rounded-[20px] p-6"
        style={{ background: "var(--surface)", border: "1px solid var(--line)", boxShadow: "var(--shadow-sm)" }}
      >
        <div className="flex items-center gap-4">
          {user.avatar_url ? (
            <img
              src={user.avatar_url}
              alt={`${user.username} 的头像`}
              className="h-16 w-16 rounded-full object-cover"
            />
          ) : (
            <div
              className="flex h-16 w-16 items-center justify-center rounded-full"
              style={{ background: "var(--surface-muted)" }}
            >
              <UserIcon className="h-8 w-8 text-ink-soft" />
            </div>
          )}
          <div>
            <h1 className="font-display text-2xl tracking-[-0.02em]">
              {user.display_name || user.username}
            </h1>
            <p className="text-sm text-ink-soft">@{user.username}</p>
            <p className="mt-1 text-xs text-ink-soft">
              加入于 {new Date(user.created_at).toLocaleDateString("zh-CN")}
            </p>
          </div>
        </div>
      </div>

      {/* Stats */}
      <div className="grid gap-4 sm:grid-cols-3">
        {[
          { icon: Lightbulb, label: "发起想法", value: profile.idea_count },
          { icon: FileText, label: "贡献方案", value: profile.contribution_count },
          { icon: Trophy, label: "精选入围", value: profile.featured_count },
        ].map((item) => (
          <div
            key={item.label}
            className="flex items-center gap-3 rounded-[14px] p-4"
            style={{ background: "var(--surface)", border: "1px solid var(--line)" }}
          >
            <div
              className="flex h-10 w-10 shrink-0 items-center justify-center rounded-[10px]"
              style={{ background: "linear-gradient(135deg, var(--accent), var(--accent-deep))" }}
            >
              <item.icon className="h-5 w-5 text-white" />
            </div>
            <div>
              <p className="font-display text-xl font-bold tracking-[-0.02em]">
                {item.value}
              </p>
              <p className="text-xs text-ink-soft">{item.label}</p>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
