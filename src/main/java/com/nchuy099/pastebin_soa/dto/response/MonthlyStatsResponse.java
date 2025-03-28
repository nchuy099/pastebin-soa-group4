package com.nchuy099.pastebin_soa.dto.response;

import com.nchuy099.pastebin_soa.dto.projection.MonthlyStatsProjection;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class MonthlyStatsResponse {
    private Long totalPastes;
    private Long totalViews;
    private Double avgViewsPerPaste;
    private Integer minViews;
    private Integer maxViews;
    private Integer activePastes;
    private Integer expiredPastes;

    public MonthlyStatsResponse(MonthlyStatsProjection projection) {
        this.totalPastes = projection.getTotalPastes() != null ? projection.getTotalPastes() : 0L;
        this.totalViews = projection.getTotalViews() != null ? projection.getTotalViews() : 0L;
        this.avgViewsPerPaste = projection.getAvgViewsPerPaste() != null ? projection.getAvgViewsPerPaste() : 0.0;
        this.minViews = projection.getMinViews() != null ? projection.getMinViews() : 0;
        this.maxViews = projection.getMaxViews() != null ? projection.getMaxViews() : 0;
        this.activePastes = projection.getActivePastes() != null ? projection.getActivePastes() : 0;
        this.expiredPastes = projection.getExpiredPastes() != null ? projection.getExpiredPastes() : 0;
    }

}
