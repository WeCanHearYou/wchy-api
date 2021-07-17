import React, { useEffect, useState } from "react"
import IconSearch from "@fider/assets/images/heroicons-search.svg"
import { Input, ShowPostStatus } from "@fider/components"
import { actions } from "@fider/services"
import { Post, PostStatus } from "@fider/models"
import { HStack, VStack } from "@fider/components/layout"
import { t, Trans } from "@lingui/macro"

interface PostSearchProps {
  exclude?: number[]
  onChanged(postNumber: number): void
}

export const PostSearch = (props: PostSearchProps) => {
  const [query, setQuery] = useState("")
  const [posts, setPosts] = useState<Post[]>([])
  const [selectedPost, setSelectedPost] = useState<Post>()

  useEffect(() => {
    if (!query) {
      return
    }

    const timer = setTimeout(() => {
      actions.searchPosts({ query, limit: 6 }).then((res) => {
        if (res.ok) {
          const filteredPosts =
            props.exclude && props.exclude.length > 0 ? res.data.filter((i) => props.exclude && props.exclude.indexOf(i.number) === -1) : res.data
          setPosts(filteredPosts)
        }
      })
    }, 500)

    return () => {
      clearTimeout(timer)
    }
  }, [query])

  const selectPost = (post: Post) => () => {
    props.onChanged(post.number)
    setSelectedPost(post)
  }

  return (
    <>
      <Input
        field="query"
        icon={IconSearch}
        placeholder={t({ id: "showpost.postsearch.query.placeholder", message: "Search original post..." })}
        value={query}
        onChange={setQuery}
      />
      <div className="grid gap-2 grid-cols-1 lg:grid-cols-3">
        {posts.map((p) => (
          <VStack onClick={selectPost(p)} className={`bg-gray-50 p-4 clickable border-2 rounded ${selectedPost === p ? "border-primary-base" : ""}`} key={p.id}>
            <HStack className="text-2xs">
              <span>#{p.number}</span> <span>&middot;</span> <ShowPostStatus status={PostStatus.Get(p.status)} /> <span>&middot;</span>{" "}
              <span>
                <Trans id="showpost.postsearch.numofvotes">{p.votesCount} votes</Trans>
              </span>
            </HStack>
            <span>{p.title}</span>
          </VStack>
        ))}
      </div>
    </>
  )
}
