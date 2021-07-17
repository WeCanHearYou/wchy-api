import React, { useState } from "react"
import { Tag } from "@fider/models"
import { Dropdown, Icon, ShowTag } from "@fider/components"
import IconCheck from "@fider/assets/images/heroicons-check.svg"
import { HStack } from "@fider/components/layout"
import { plural, t, Trans } from "@lingui/macro"

interface TagsFilterProps {
  tags: Tag[]
  selected: string[]
  selectionChanged: (selected: string[]) => void
}

export const TagsFilter = (props: TagsFilterProps) => {
  const [selected, setSelected] = useState<string[]>(props.selected)

  const toggle = (item: Tag) => () => {
    const idx = selected.indexOf(item.slug)
    const next = idx >= 0 ? selected.splice(idx, 1) && selected : selected.concat(item.slug)
    setSelected(next)
    props.selectionChanged(next)
  }

  const count = selected.length
  const label =
    count >= 1
      ? plural(count, {
          one: `# tag`,
          other: `# tags`,
        })
      : t({ id: "home.tagsfilter.selected.none", message: "Any tag" })

  return (
    <HStack>
      <Trans id="home.tagsfilter.label.with">with</Trans>
      <Dropdown renderHandle={<div className="text-medium text-xs uppercase rounded-md uppercase bg-gray-100 px-2 py-1">{label}</div>}>
        {props.tags.map((t) => (
          <Dropdown.ListItem onClick={toggle(t)} key={t.id}>
            <HStack spacing={2}>
              <Icon sprite={IconCheck} className={`h-4 text-green-600 ${!selected.includes(t.slug) && "invisible"}`} />
              <ShowTag tag={t} circular={true} />
              <span>{t.name}</span>
            </HStack>
          </Dropdown.ListItem>
        ))}
      </Dropdown>
    </HStack>
  )
}
