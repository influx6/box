// Package docker implements a higher level API ontop of the moby client package
// to provide easy interactions on docker images and containers.
/*

@templater(id => Spell, gen => Partial.Go, file => _spell.tml)

@templaterTypesFor(asJSON, id => Spell, filename => checkpoint_create.go, Name => CheckpointCreate, {
   {
       "return": [],
       "arguments": ["container string", "chop types.CheckpointCreateOptions"]
   }
})

@templaterTypesFor(asJSON, id => Spell, filename => checkpoint_delete.go, Name => CheckpointDelete, {
   {
       "return": [],
       "arguments": ["container string", "chop types.CheckpointDeleteOptions"]
   }
})

@templaterTypesFor(asJSON, id => Spell, filename => image_tag.go, Name => ImageTag, {
   {
       "return": [],
       "arguments": ["tag string"]
   }
})

@templaterTypesFor(asJSON, id => Spell, filename => create_image.go, Name => CreateImage, {
   {
       "return": ["types.ImageLoadResponse"],
       "arguments": ["reader io.ReadCloser"]
   }
})

@templaterTypesFor(asJSON, id => Spell, filename => image_history.go, Name => ImageHistory,  {
   {
       "return": ["image.HistoryResponseItem"],
       "arguments": []
   }
})

@templaterTypesFor(asJSON, id => Spell, filename => image_inspect_with_raw.go, Name => ImageInspectWithRaw, {
   {
       "return": ["types.ImageInspect"],
       "arguments": []
   }
})

@templaterTypesFor(asJSON, id => Spell, filename => events.go, Name => Events, OptionType => types.EventsOptions, {
   {
       "return": [],
       "arguments": ["eventOp types.EventsOptions"]
   }
})

@templaterTypesFor(asJSON, id => Spell, filename => image_save.go, Name => ImageSave, {
   {
       "return": ["io.ReadCloser"],
       "arguments": ["ops []string"]
   }
})

@templaterTypesFor(asJSON, id => Spell, filename => image_pull.go, Name => ImagePull, {
   {
       "return": ["io.ReadCloser"],
       "arguments": ["imgOp types.ImagePullOptions"]
   }
})

@templaterTypesFor(asJSON, id => Spell, filename => image_push.go, Name => ImagePush, {
   {
       "return": ["io.ReadCloser"],
       "arguments": ["imp types.ImagePushOptions"]
   }
})

@templaterTypesFor(asJSON, id => Spell, filename => image_prune.go, Name => ImagePrune, {
   {
       "return": ["types.ImagesPruneReport"],
       "arguments": ["args filters.Args"]
   }
})

@templaterTypesFor(asJSON, id => Spell, filename => image_load.go, Name => ImageLoad, {
   {
       "return": ["types.ImageLoadResponse"],
       "arguments": ["reader io.Reader"]
   }
})

@templaterTypesFor(asJSON, id => Spell, filename => image_import.go, Name => ImageImport, {
   {
       "return": ["io.ReadCloser"],
       "arguments": ["impOp types.ImageImportOptions"]
   }
})

@templaterTypesFor(asJSON, id => Spell, filename => image_list.go, Name => ImageList, {
   {
       "return": ["[]types.ImageSummary"],
       "arguments": ["listOps types.ImageListOptions"]
   }
})

@templaterTypesFor(asJSON, id => Spell, filename => image_search.go, Name => ImageSearch, {
   {
       "return": ["[]registry.SearchResult"],
       "arguments": ["searchOps types.ImageSearchOptions"]
   }
})

@templaterTypesFor(asJSON, id => Spell, asJSON, filename => image_remove.go, Name => ImageRemove, {
   {
       "return": ["[]types.ImageDeleteResponseItem"],
       "arguments": ["removeOps types.ImageRemoveOptions"]
   }
})

@templaterTypesFor(asJSON, id => Spell, asJSON, filename => container_wait.go, Name => ContainerWait, {
   {
       "arguments": ["containerID string", "container container.WaitCondition"],
       "returns": ["<-chan container.ContainerWaitOKBody", "<-chan error"]
   }
})

@templaterTypesFor(asJSON, id => Spell, asJSON, filename => container_copy_to.go, Name => CopyToContainer, {
   {
       "arguments": ["container string", "topath string", "reader io.ReadCloser", "cops types.CopyToContainerOptions"],
       "return": []
   }
 })

@templaterTypesFor(asJSON, id => Spell, asJSON, filename => container_copy_from.go, Name => CopyFromContainer, {
   {
       "arguments": ["container string", "srcPath string"],
       "return": ["io.ReadCloser", "types.ContainerPathStat"]
   }
 })

@templaterTypesFor(asJSON, id => Spell, filename => network_disconnect.go, Name => NetworkDisconnect,  {
  {
       "arguments": ["networkID string"],
       "return": []
  }
})

@templaterTypesFor(asJSON, id => Spell, filename => network_create.go, Name => NetworkCreate, {
  {
     "arguments": ["network types.NetworkCreate"],
     "return": ["types.NetworkCreateResponse"]
  }
})

@templaterTypesFor(asJSON, id => Spell, filename => network_inspect.go, Name => NetworkInspect, {
  {
     "arguments": ["netOp types.NetworkInspectOptions"],
     "return": ["types.NetworkResource"]
  }
})

*
*/
package docker
