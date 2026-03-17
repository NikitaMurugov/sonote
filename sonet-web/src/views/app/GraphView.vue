<template>
  <div class="h-full flex flex-col anim-fade-up">
    <div class="px-6 py-4 flex items-center justify-between border-b border-border-light/60">
      <h1 class="text-lg font-semibold" style="font-family: var(--font-heading)">Граф связей</h1>
      <span class="text-[11px] text-text-tertiary">{{ nodes.length }} заметок, {{ edges.length }} связей</span>
    </div>
    <div ref="graphContainer" class="flex-1 relative overflow-hidden bg-bg-base">
      <svg ref="svgEl" class="w-full h-full"></svg>
      <!-- Tooltip -->
      <div
        v-if="hoveredNode"
        class="absolute pointer-events-none bg-bg-surface border border-border rounded-lg px-3 py-1.5 text-sm text-text-primary shadow-lg"
        :style="{ left: tooltipX + 'px', top: tooltipY + 'px' }"
      >
        {{ hoveredNode.title }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useWorkspaceStore } from '@/stores/workspace'
import api from '@/composables/useApi'
import { forceSimulation, forceLink, forceManyBody, forceCenter, forceCollide } from 'd3-force'
import { select } from 'd3-selection'
import { zoom } from 'd3-zoom'

interface GraphNode {
  id: number
  title: string
  slug: string
  x?: number
  y?: number
  vx?: number
  vy?: number
}

interface GraphEdge {
  source: number | GraphNode
  target: number | GraphNode
}

const route = useRoute()
const router = useRouter()
const workspaceStore = useWorkspaceStore()

const graphContainer = ref<HTMLElement>()
const svgEl = ref<SVGSVGElement>()
const nodes = ref<GraphNode[]>([])
const edges = ref<GraphEdge[]>([])
const hoveredNode = ref<GraphNode | null>(null)
const tooltipX = ref(0)
const tooltipY = ref(0)

let simulation: any = null

async function fetchGraph() {
  if (!workspaceStore.currentWorkspace) return
  const { data } = await api.get(`/workspaces/${workspaceStore.currentWorkspace.id}/graph`)
  nodes.value = data.data.nodes || []
  edges.value = data.data.edges || []
  renderGraph()
}

function renderGraph() {
  if (!svgEl.value || !nodes.value.length) return

  const svg = select(svgEl.value)
  svg.selectAll('*').remove()

  const width = svgEl.value.clientWidth
  const height = svgEl.value.clientHeight

  const g = svg.append('g')

  // Zoom
  const zoomBehavior = zoom<SVGSVGElement, unknown>()
    .scaleExtent([0.1, 4])
    .on('zoom', (event) => g.attr('transform', event.transform))
  svg.call(zoomBehavior)

  // Get CSS vars
  const style = getComputedStyle(document.documentElement)
  const colorAccent = style.getPropertyValue('--color-accent-soft').trim()
  const colorPrimary = style.getPropertyValue('--color-primary').trim()
  const colorBorder = style.getPropertyValue('--color-border').trim()
  const colorText = style.getPropertyValue('--color-text-secondary').trim()

  // Count connections per node
  const connectionCount = new Map<number, number>()
  edges.value.forEach((e) => {
    const s = typeof e.source === 'object' ? e.source.id : e.source
    const t = typeof e.target === 'object' ? e.target.id : e.target
    connectionCount.set(s, (connectionCount.get(s) || 0) + 1)
    connectionCount.set(t, (connectionCount.get(t) || 0) + 1)
  })

  simulation = forceSimulation(nodes.value as any)
    .force('link', forceLink(edges.value as any).id((d: any) => d.id).distance(80))
    .force('charge', forceManyBody().strength(-200))
    .force('center', forceCenter(width / 2, height / 2))
    .force('collide', forceCollide().radius(30))

  const link = g.append('g')
    .selectAll('line')
    .data(edges.value)
    .join('line')
    .attr('stroke', colorBorder)
    .attr('stroke-width', 1)
    .attr('stroke-opacity', 0.5)

  const node = g.append('g')
    .selectAll('circle')
    .data(nodes.value)
    .join('circle')
    .attr('r', (d: GraphNode) => Math.min(6 + (connectionCount.get(d.id) || 0) * 2, 16))
    .attr('fill', colorAccent)
    .attr('stroke', colorBorder)
    .attr('stroke-width', 1)
    .attr('cursor', 'pointer')
    .on('mouseover', (_: any, d: GraphNode) => {
      hoveredNode.value = d
      select(_ .target).attr('fill', colorPrimary).attr('r', Math.min(8 + (connectionCount.get(d.id) || 0) * 2, 18))
    })
    .on('mousemove', (event: MouseEvent) => {
      tooltipX.value = event.offsetX + 12
      tooltipY.value = event.offsetY - 10
    })
    .on('mouseout', (_: any, d: GraphNode) => {
      hoveredNode.value = null
      select(_.target).attr('fill', colorAccent).attr('r', Math.min(6 + (connectionCount.get(d.id) || 0) * 2, 16))
    })
    .on('click', (_: any, d: GraphNode) => {
      router.push(`/w/${route.params.wsSlug}/note/${d.id}`)
    })

  const label = g.append('g')
    .selectAll('text')
    .data(nodes.value)
    .join('text')
    .text((d: GraphNode) => d.title.length > 20 ? d.title.slice(0, 20) + '…' : d.title)
    .attr('font-size', '10px')
    .attr('fill', colorText)
    .attr('text-anchor', 'middle')
    .attr('dy', (d: GraphNode) => -(Math.min(6 + (connectionCount.get(d.id) || 0) * 2, 16)) - 4)
    .attr('pointer-events', 'none')
    .style('font-family', 'var(--font-body)')

  simulation.on('tick', () => {
    link
      .attr('x1', (d: any) => d.source.x)
      .attr('y1', (d: any) => d.source.y)
      .attr('x2', (d: any) => d.target.x)
      .attr('y2', (d: any) => d.target.y)
    node.attr('cx', (d: any) => d.x).attr('cy', (d: any) => d.y)
    label.attr('x', (d: any) => d.x).attr('y', (d: any) => d.y)
  })
}

onMounted(async () => {
  const slug = route.params.wsSlug as string
  await workspaceStore.setCurrentBySlug(slug)
  await fetchGraph()
})

onUnmounted(() => {
  simulation?.stop()
})
</script>
