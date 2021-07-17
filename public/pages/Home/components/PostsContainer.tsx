import "./PostsContainer.scss"

import React from "react"

import { Post, Tag, CurrentUser } from "@fider/models"
import { Loader, Input } from "@fider/components"
import { actions, navigator, querystring } from "@fider/services"
import IconSearch from "@fider/assets/images/heroicons-search.svg"
import IconX from "@fider/assets/images/heroicons-x.svg"
import { PostFilter } from "./PostFilter"
import { ListPosts } from "./ListPosts"
import { TagsFilter } from "./TagsFilter"
import { t, Trans } from "@lingui/macro"

interface PostsContainerProps {
  user?: CurrentUser
  posts: Post[]
  tags: Tag[]
  countPerStatus: { [key: string]: number }
}

interface PostsContainerState {
  loading: boolean
  posts?: Post[]
  view: string
  tags: string[]
  query: string
  limit?: number
}

export class PostsContainer extends React.Component<PostsContainerProps, PostsContainerState> {
  constructor(props: PostsContainerProps) {
    super(props)

    this.state = {
      posts: this.props.posts,
      loading: false,
      view: querystring.get("view"),
      query: querystring.get("query"),
      tags: querystring.getArray("tags"),
      limit: querystring.getNumber("limit"),
    }
  }

  private changeFilterCriteria<K extends keyof PostsContainerState>(obj: Pick<PostsContainerState, K>, reset: boolean): void {
    this.setState(obj, () => {
      const query = this.state.query.trim().toLowerCase()
      navigator.replaceState(
        querystring.stringify({
          tags: this.state.tags,
          query,
          view: this.state.view,
          limit: this.state.limit,
        })
      )

      this.searchPosts(query, this.state.view, this.state.limit, this.state.tags, reset)
    })
  }

  private timer?: number
  private async searchPosts(query: string, view: string, limit: number | undefined, tags: string[], reset: boolean) {
    window.clearTimeout(this.timer)
    this.setState({ posts: reset ? undefined : this.state.posts, loading: true })
    this.timer = window.setTimeout(() => {
      actions.searchPosts({ query, view, limit, tags }).then((response) => {
        if (response.ok && this.state.loading) {
          this.setState({ loading: false, posts: response.data })
        }
      })
    }, 500)
  }

  private handleViewChanged = (view: string) => {
    this.changeFilterCriteria({ view }, true)
  }

  private handleTagsFilterChanged = (tags: string[]) => {
    this.changeFilterCriteria({ tags }, true)
  }

  private handleSearchFilterChanged = (query: string) => {
    this.changeFilterCriteria({ query }, true)
  }

  private clearSearch = () => {
    this.changeFilterCriteria({ query: "" }, true)
  }

  private showMore = (event: React.MouseEvent<HTMLElement> | React.TouchEvent<HTMLElement>): void => {
    event.preventDefault()
    this.changeFilterCriteria({ limit: (this.state.limit || 30) + 10 }, false)
  }

  private getShowMoreLink = (): string | undefined => {
    if (this.state.posts && this.state.posts.length >= (this.state.limit || 30)) {
      return querystring.set("limit", (this.state.limit || 30) + 10)
    }
  }

  public render() {
    const showMoreLink = this.getShowMoreLink()

    return (
      <div className="c-posts-container">
        <div className="c-posts-container__header mb-4">
          {!this.state.query && (
            <div className="c-posts-container__filter-col">
              <PostFilter activeView={this.state.view} viewChanged={this.handleViewChanged} countPerStatus={this.props.countPerStatus} />
              {this.props.tags.length > 0 && <TagsFilter tags={this.props.tags} selectionChanged={this.handleTagsFilterChanged} selected={this.state.tags} />}
            </div>
          )}
          <div className="c-posts-container__search-col">
            <Input
              field="query"
              icon={this.state.query ? IconX : IconSearch}
              onIconClick={this.state.query ? this.clearSearch : undefined}
              placeholder={t({ id: "home.postscontainer.query.placeholder", message: "Search" })}
              value={this.state.query}
              onChange={this.handleSearchFilterChanged}
            />
          </div>
        </div>
        <ListPosts
          posts={this.state.posts}
          tags={this.props.tags}
          emptyText={t({ id: "home.postscontainer.label.noresults", message: "No results matched your search, try something different." })}
        />
        {this.state.loading && <Loader />}
        {showMoreLink && (
          <div className="my-4 ml-4">
            <a href={showMoreLink} className="text-primary-base text-medium hover:underline" onClick={this.showMore}>
              <Trans id="home.postscontainer.label.viewmore">View more posts</Trans>
            </a>
          </div>
        )}
      </div>
    )
  }
}
